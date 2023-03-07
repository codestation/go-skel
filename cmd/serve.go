// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.dev/go/go-skel/app"
	"megpoid.dev/go/go-skel/config"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start service",
	Long:  `Starts the HTTP endpoint and other services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	// show version on console
	printVersion()

	// load config
	cfg, err := config.NewConfig(config.WithUnmarshal(unmarshalFunc))
	if err != nil {
		return err
	}

	// setup channel to check when app is stopped
	quit := make(chan os.Signal, 1)

	// Create a new newApp
	newApp, err := app.NewApp(cfg)
	if err != nil {
		return fmt.Errorf("cannot create newApp: %w", err)
	}
	defer newApp.Shutdown()

	// Start newApp
	if err := newApp.Start(); err != nil {
		return fmt.Errorf("cannot start newApp: %w", err)
	}

	// Wait for kill signal before attempting to gracefully stop the running service
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	generalFs := config.LoadGeneralFlags()
	serverFs := config.LoadServerFlags()
	databaseFs := config.LoadDatabaseFlags()

	serveCmd.Flags().AddFlagSet(generalFs)
	serveCmd.Flags().AddFlagSet(serverFs)
	serveCmd.Flags().AddFlagSet(databaseFs)

	cobra.CheckErr(viper.BindPFlags(generalFs))
	cobra.CheckErr(viper.BindPFlags(serverFs))
	cobra.CheckErr(viper.BindPFlags(databaseFs))
}

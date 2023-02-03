// Copyright 2022 codestation. All rights reserved.
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
	"megpoid.dev/go/go-skel/api"
	"megpoid.dev/go/go-skel/app"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/web"
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

	// Create a new server
	server, err := app.NewServer(cfg)
	if err != nil {
		return fmt.Errorf("cannot create server: %w", err)
	}
	defer server.Shutdown()

	// Create web, websocket, graphql, etc, servers
	_, err = api.Init(server)
	if err != nil {
		return fmt.Errorf("cannot initialize API: %w", err)
	}
	web.New(server)

	// Start server
	if err := server.Start(); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	// Wait for kill signal before attempting to gracefully stop the running service
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().String("jwt-secret", "", "JWT secret key")
	serveCmd.Flags().String("encryption-key", "", "Application encryption key")
	serveCmd.Flags().StringP("listen", "l", config.DefaultListenAddress, "Listen address")
	serveCmd.Flags().DurationP("timeout", "t", config.DefaultReadTimeout, "Request timeout")
	serveCmd.Flags().String("dsn", "", "Database connection string. Setting the DSN ignores the db-* settings")
	serveCmd.Flags().Int("query-limit", 1000, "Max results per query")
	serveCmd.Flags().Int("max-open-conns", config.DefaultMaxOpenConns, "Max open connections")
	serveCmd.Flags().Int("max-idle-conns", config.DefaultMaxIdleConns, "Max idle connections")
	serveCmd.Flags().String("body-limit", "", "Max body size for http requests")
	serveCmd.Flags().StringSlice("cors-allow-origin", []string{}, "CORS Allowed origins")
	err := viper.BindPFlags(serveCmd.Flags())
	cobra.CheckErr(err)
}

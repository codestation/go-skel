/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/api"
	"megpoid.xyz/go/go-skel/app"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/web"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start service",
	Long:  `Starts the HTTP endpoint and other services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		quit := make(chan os.Signal, 1)

		// load config
		var cfg model.Config
		if err := viper.Unmarshal(&cfg.GeneralSettings, unmarshalDecoders...); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.ServerSettings, unmarshalDecoders...); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.SqlSettings, unmarshalDecoders...); err != nil {
			return err
		}
		cfg.SetDefaults()
		if err := cfg.Validate(); err != nil {
			return err
		}
		printVersion()

		return runServer(cfg, quit)
	},
}

func runServer(cfg model.Config, quit chan os.Signal) error {
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
	serveCmd.Flags().StringP("listen", "l", ":8000", "Listen address")
	serveCmd.Flags().DurationP("timeout", "t", 30*time.Second, "Request timeout")
	serveCmd.Flags().String("dsn", "", "Database connection string. Setting the DSN ignores the db-* settings")
	serveCmd.Flags().Int("query-limit", 1000, "Max results per query")
	serveCmd.Flags().StringSlice("cors-allow-origin", []string{}, "CORS Allowed origins")
	err := viper.BindPFlags(serveCmd.Flags())
	cobra.CheckErr(err)
}

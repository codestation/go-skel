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

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/api"
	"megpoid.xyz/go/go-skel/app"
	"megpoid.xyz/go/go-skel/config"
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
		var cfg config.Config
		if err := viper.Unmarshal(&cfg, viper.DecodeHook(HexStringToByteArray())); err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		printVersion()

		return runServer(&cfg, quit)
	},
}

func runServer(cfg *config.Config, quit chan os.Signal) error {
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

	serveCmd.Flags().StringP("jwt-secret", "", "", "JWT secret key")
	serveCmd.Flags().StringP("master-key", "", "", "Application master key")
	serveCmd.Flags().StringP("listen", "l", ":8000", "Listen address")
	serveCmd.Flags().DurationP("timeout", "t", 30*time.Second, "Request timeout")
	serveCmd.Flags().StringP("dsn", "n", "", "Database connection string. Setting the DSN ignores the db-* settings")
	serveCmd.Flags().StringP("db-adapter", "a", "postgres", "Database adapter")
	serveCmd.Flags().StringP("db-host", "H", "localhost", "Database host")
	serveCmd.Flags().StringP("db-port", "p", "5432", "Database port")
	serveCmd.Flags().StringP("db-name", "m", "", "Database name")
	serveCmd.Flags().StringP("db-user", "u", "postgres", "Database user")
	serveCmd.Flags().StringP("db-password", "w", "", "Database password")
	err := viper.BindPFlags(serveCmd.Flags())
	cobra.CheckErr(err)
}

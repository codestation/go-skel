// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/api"
	"megpoid.xyz/go/go-skel/app"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/web"
	"os"
	"os/signal"
	"syscall"
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
	var cfg model.Config
	if err := viper.Unmarshal(&cfg.GeneralSettings, unmarshalDecoder); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg.ServerSettings, unmarshalDecoder); err != nil {
		return err
	}
	if err := viper.Unmarshal(&cfg.SqlSettings, unmarshalDecoder); err != nil {
		return err
	}
	cfg.SetDefaults()
	if err := cfg.Validate(); err != nil {
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
	serveCmd.Flags().StringP("listen", "l", model.DefaultListenAddress, "Listen address")
	serveCmd.Flags().DurationP("timeout", "t", model.DefaultReadTimeout, "Request timeout")
	serveCmd.Flags().String("dsn", "", "Database connection string. Setting the DSN ignores the db-* settings")
	serveCmd.Flags().Int("query-limit", 1000, "Max results per query")
	serveCmd.Flags().Int("max-open-conns", model.DefaultMaxOpenConns, "Max open connections")
	serveCmd.Flags().Int("max-idle-conns", model.DefaultMaxIdleConns, "Max idle connections")
	serveCmd.Flags().String("body-limit", "", "Max body size for http requests")
	serveCmd.Flags().StringSlice("cors-allow-origin", []string{}, "CORS Allowed origins")
	err := viper.BindPFlags(serveCmd.Flags())
	cobra.CheckErr(err)
}

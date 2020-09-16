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
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"log"
	"megpoid.xyz/go/go-skel/internal"
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/services"
	"megpoid.xyz/go/go-skel/pkg/hooks"
	"megpoid.xyz/go/go-skel/pkg/sql"
	"megpoid.xyz/go/go-skel/pkg/sql/connection"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var defaultDSN = "host=127.0.0.1 port=5432 user=postgres password=secret dbname=app"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var cfg config.Config
		err := viper.Unmarshal(&cfg, viper.DecodeHook(hooks.HexStringToByteArray()))
		if err != nil {
			return err
		}

		if cfg.MasterKey == nil {
			return errors.New("no application master key specified")
		}
		if cfg.JWTSecret == nil {
			return errors.New("no jwt secret key specified")
		}
		if cfg.DSN == "" {
			return errors.New("no dsn specified")
		}

		db, err := sql.NewDatabase(ctx, cfg.DSN, sql.Config{MaxOpenConnections: 5})
		if err != nil {
			return err
		}
		defer func() {
			err := db.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		e := echo.New()
		e.HideBanner = true
		e.Debug = cfg.Debug
		e.Use(middleware.Gzip())
		e.Use(middleware.Recover())

		if e.Debug {
			e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format: "method=${method}, uri=${uri}, status=${status}\n",
			}))
		}

		svr := services.NewHTTPServer(cfg.Addr, e)
		module := internal.New(connection.New(db), &cfg)
		module.Handler.Register(e)

		go func() {
			if err := svr.Serve(); err != nil {
				log.Fatal(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		if err := svr.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("jwt-secret", "", "", "JWT secret key")
	serveCmd.Flags().StringP("master-key", "", "", "Application master key")
	serveCmd.Flags().StringP("listen", "l", ":8000", "Listen address")
	serveCmd.Flags().StringP("dsn", "n", defaultDSN, "Database connection string")
	err := viper.BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}

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
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/internal"
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/services"
	"megpoid.xyz/go/go-skel/pkg/hooks"
	"megpoid.xyz/go/go-skel/pkg/sql"
	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start service",
	Long:  `Starts the HTTP endpoint and other services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printVersion()

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

		db, err := sql.NewDatabase(ctx, cfg.DBAdapter, cfg.GetDSN(), sql.MaxOpenConns(5))
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
		module := internal.New(connection.NewPostgresConn(db), &cfg)
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

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/db"
	"megpoid.dev/go/go-skel/pkg/sql"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(config.WithUnmarshal(unmarshalFunc))
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		quit := make(chan os.Signal, 1)

		// Database initialization
		pool, err := sql.NewConnection(sql.Config{
			DataSourceName:  cfg.DatabaseSettings.DataSourceName,
			MaxIdleConns:    cfg.DatabaseSettings.MaxIdleConns,
			MaxOpenConns:    cfg.DatabaseSettings.MaxOpenConns,
			ConnMaxLifetime: cfg.DatabaseSettings.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.DatabaseSettings.ConnMaxIdleTime,
			QueryLimit:      cfg.DatabaseSettings.QueryLimit,
		})
		if err != nil {
			return err
		}
		defer pool.Close()

		go func() {
			err := db.RunMigrations(ctx, pool, cfg)
			if err != nil {
				log.Println(err.Error())
			}

			quit <- os.Interrupt
		}()

		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		cancel()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	databaseFlags := config.LoadDatabaseFlags()
	migrateFlags := config.LoadMigrateFlags()

	migrateCmd.Flags().AddFlagSet(databaseFlags)
	migrateCmd.Flags().AddFlagSet(migrateFlags)

	cobra.CheckErr(viper.BindPFlags(databaseFlags))
	cobra.CheckErr(viper.BindPFlags(migrateFlags))
}

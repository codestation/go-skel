// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/db"
	"megpoid.dev/go/go-skel/pkg/cfg"
	"megpoid.dev/go/go-skel/pkg/migration"
	"megpoid.dev/go/go-skel/pkg/sql"
	"megpoid.dev/go/go-skel/testdata"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var slogOpts *slog.HandlerOptions
		if viper.GetBool("debug") {
			slogOpts = &slog.HandlerOptions{Level: slog.LevelDebug}
		}

		handler := slog.NewTextHandler(os.Stdout, slogOpts)
		slog.SetDefault(slog.New(handler))

		databaseSettings := config.DatabaseSettings{}
		if err := cfg.ReadConfig(&databaseSettings); err != nil {
			return fmt.Errorf("failed to read database settings: %w", err)
		}

		migrationSettings := config.MigrationSettings{}
		if err := cfg.ReadConfig(&migrationSettings); err != nil {
			return fmt.Errorf("failed to read migration settings: %w", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		quit := make(chan os.Signal, 1)

		// Database initialization
		pool, err := sql.NewConnection(sql.Config{
			DataSourceName:  databaseSettings.DataSourceName,
			MaxIdleConns:    databaseSettings.MaxIdleConns,
			MaxOpenConns:    databaseSettings.MaxOpenConns,
			ConnMaxLifetime: databaseSettings.ConnMaxLifetime,
			ConnMaxIdleTime: databaseSettings.ConnMaxIdleTime,
			QueryLimit:      databaseSettings.QueryLimit,
		})
		if err != nil {
			return err
		}
		defer pool.Close()

		var migrationErr error

		migrationConfig := migration.Options{
			TableName: "app_migrations",
			Redo:      migrationSettings.Redo,
			Reset:     migrationSettings.Reset,
			Rollback:  migrationSettings.Rollback,
			Seed:      migrationSettings.Seed,
			Step:      migrationSettings.Step,
			Test:      migrationSettings.Test,
			MigrationAsset: migration.AssetOptions{
				FS:   db.Assets(),
				Root: "migrations",
			},
			SeedAsset: migration.AssetOptions{
				FS:   db.Seeds(),
				Root: "seed",
			},
			TestAsset: migration.AssetOptions{
				FS:   testdata.SqlAssets(),
				Root: "sql",
			},
		}

		go func() {
			migrationErr = migration.RunMigrations(ctx, pool, migrationConfig)
			if migrationErr != nil {
				slog.Error(migrationErr.Error())
			}

			quit <- os.Interrupt
		}()

		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		cancel()

		return migrationErr
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	databaseFlags := config.LoadDatabaseFlags(migrateCmd.Name())
	migrateFlags := config.LoadMigrateFlags(migrateCmd.Name())

	migrateCmd.Flags().AddFlagSet(databaseFlags)
	migrateCmd.Flags().AddFlagSet(migrateFlags)
}

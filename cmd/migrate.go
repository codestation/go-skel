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
	"log"
	"megpoid.xyz/go/go-skel/config"

	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/db"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var cfg config.Config
		err := viper.Unmarshal(&cfg, viper.DecodeHook(HexStringToByteArray()))
		if err != nil {
			return err
		}

		migrations := migrate.EmbedFileSystemMigrationSource{
			FileSystem: db.Assets(),
			Root:       "migrations",
		}

		db, err := sqlx.ConnectContext(ctx, cfg.DBAdapter, cfg.GetDSN())
		if err != nil {
			return err
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Print(err)
			}
		}()

		migrate.SetTable("app_migrations")

		if viper.GetBool("reset") {
			_, err := db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
			if err != nil {
				return err
			}

			_, err = db.ExecContext(ctx, "CREATE SCHEMA public")
			if err != nil {
				return nil
			}
			log.Printf("Recreated 'public' schema")
		}

		step := 0

		if !viper.GetBool("reset") && (viper.GetBool("rollback") || viper.GetBool("redo")) {
			step = viper.GetInt("step")
			n, err := migrate.ExecMax(db.DB, cfg.DBAdapter, migrations, migrate.Down, step)
			if err != nil {
				return err
			}
			log.Printf("Reverted %d migrations", n)
		}

		if viper.GetBool("reset") || !viper.GetBool("rollback") || viper.GetBool("redo") {
			if viper.GetBool("redo") {
				step = viper.GetInt("step")
			}

			n, err := migrate.ExecMax(db.DB, cfg.DBAdapter, migrations, migrate.Up, step)
			if err != nil {
				return err
			}
			log.Printf("Applied %d migrations", n)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("rollback", "r", false, "Rollback last migration")
	migrateCmd.Flags().BoolP("redo", "e", false, "Rollback last migration then migrate again")
	migrateCmd.Flags().BoolP("reset", "t", false, "Drop all tables and run migration")
	migrateCmd.Flags().IntP("step", "s", 1, "Steps to rollback/redo")
	migrateCmd.Flags().StringP("dsn", "n", "", "Database connection string. Setting the DSN ignores the db-* settings")
	migrateCmd.Flags().StringP("db-adapter", "a", "postgres", "Database adapter")
	migrateCmd.Flags().StringP("db-host", "h", "localhost", "Database host")
	migrateCmd.Flags().StringP("db-port", "p", "5432", "Database port")
	migrateCmd.Flags().StringP("db-name", "m", "", "Database name")
	migrateCmd.Flags().StringP("db-user", "u", "postgres", "Database user")
	migrateCmd.Flags().StringP("db-password", "w", "", "Database password")
	err := viper.BindPFlags(migrateCmd.Flags())
	cobra.CheckErr(err)
}

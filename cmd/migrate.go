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
	"log"

	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/pkg/hooks"

	"github.com/gobuffalo/packr/v2"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg config.Config
		err := viper.Unmarshal(&cfg, viper.DecodeHook(hooks.HexStringToByteArray()))
		if err != nil {
			return err
		}

		migrations := migrate.PackrMigrationSource{
			Box: packr.New("migrations", "../migrations"),
		}

		db, err := sqlx.Connect("postgres", cfg.GetDSN())
		if err != nil {
			return err
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Print(err)
			}
		}()

		migrate.SetTable("govote_migrations")

		if viper.GetBool("rollback") || viper.GetBool("redo") {
			n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Down)
			if err != nil {
				return err
			}
			log.Printf("Reverted %d migrations", n)
		}

		if !viper.GetBool("rollback") || viper.GetBool("redo") {
			n, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
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
	migrateCmd.Flags().StringP("dsn", "n", "", "Database connection string. Setting the DSN ignores the db-* settings")
	migrateCmd.Flags().StringP("db-host", "h", "localhost", "Database host")
	migrateCmd.Flags().StringP("db-port", "p", "5432", "Database port")
	migrateCmd.Flags().StringP("db-name", "m", "", "Database name")
	migrateCmd.Flags().StringP("db-user", "u", "postgres", "Database user")
	migrateCmd.Flags().StringP("db-password", "w", "", "Database password")
	err := viper.BindPFlags(migrateCmd.Flags())
	if err != nil {
		panic(err)
	}
}

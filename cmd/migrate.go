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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store/sqlstore"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  `Apply the database migrations to the database`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg model.Config
		if err := viper.Unmarshal(&cfg.GeneralSettings, unmarshalDecoders...); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.SqlSettings, unmarshalDecoders...); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.MigrationSettings, unmarshalDecoders...); err != nil {
			return err
		}
		cfg.SetDefaults()
		if err := cfg.Validate(); err != nil {
			return err
		}

		store := sqlstore.New(cfg.SqlSettings)
		defer func(store *sqlstore.SqlStore) {
			err := store.Close()
			if err != nil {
				log.Print(err)
			}
		}(store)

		return store.RunMigrations(cfg.MigrationSettings)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().Bool("rollback", false, "Rollback last migration")
	migrateCmd.Flags().Bool("redo", false, "Rollback last migration then migrate again")
	migrateCmd.Flags().Bool("reset", false, "Drop all tables and run migration")
	migrateCmd.Flags().Int("step", 1, "Steps to rollback/redo")
	migrateCmd.Flags().String("dsn", "", "Database connection string. Setting the DSN ignores the db-* settings")
	migrateCmd.Flags().String("driver", "postgres", "Database driver")
	err := viper.BindPFlags(migrateCmd.Flags())
	cobra.CheckErr(err)
}

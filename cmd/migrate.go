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
		if err := viper.Unmarshal(&cfg.GeneralSettings, unmarshalDecoder); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.SqlSettings, unmarshalDecoder); err != nil {
			return err
		}
		if err := viper.Unmarshal(&cfg.MigrationSettings, unmarshalDecoder); err != nil {
			return err
		}
		cfg.SetDefaults()
		if err := cfg.Validate(); err != nil {
			return err
		}

		// Database initialization
		conn := sqlstore.NewConnection(cfg.SqlSettings)
		defer conn.Close()

		return sqlstore.RunMigrations(conn, cfg)
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

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.dev/go/go-skel/app/tasks"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/pkg/sql"
)

// migrateCmd represents the migrate command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Run task queue",
	Long:  `Run the task queue`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(config.WithUnmarshal(unmarshalFunc))
		if err != nil {
			return err
		}

		// Database initialization
		pool, err := sql.NewConnection(cfg.DatabaseSettings)
		if err != nil {
			return err
		}
		defer pool.Close()

		conn := sql.NewPgxWrapper(pool)

		queue := asynq.NewServer(
			asynq.RedisClientOpt{Addr: cfg.GeneralSettings.RedisAddr},
			asynq.Config{
				Concurrency: cfg.GeneralSettings.Workers,
			},
		)

		mux := asynq.NewServeMux()
		mux.HandleFunc(tasks.TypeSayHello, tasks.HandleSayHelloTask)
		mux.Handle(tasks.TypeProfileCheck, tasks.NewProfileProcessor(conn))

		if err := queue.Run(mux); err != nil {
			if err != asynq.ErrServerClosed {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)

	databaseFlags := config.LoadDatabaseFlags()
	generalFlags := config.LoadGeneralFlags()

	queueCmd.Flags().AddFlagSet(databaseFlags)
	queueCmd.Flags().AddFlagSet(generalFlags)

	cobra.CheckErr(viper.BindPFlags(databaseFlags))
	cobra.CheckErr(viper.BindPFlags(generalFlags))
}

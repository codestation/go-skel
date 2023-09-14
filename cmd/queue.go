// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"errors"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"megpoid.dev/go/go-skel/app/tasks"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
)

// migrateCmd represents the migrate command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "Run task queue",
	Long:  `Run the task queue`,
	PreRun: func(cmd *cobra.Command, _ []string) {
		cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.NewConfig(config.WithUnmarshal(unmarshalFunc))
		if err != nil {
			return err
		}

		queue := asynq.NewServer(
			asynq.RedisClientOpt{Addr: cfg.GeneralSettings.RedisAddr},
			asynq.Config{
				Concurrency: cfg.GeneralSettings.Workers,
			},
		)

		backgroundUsecase := usecase.NewDelay()

		mux := asynq.NewServeMux()
		mux.Handle(tasks.TypeDelay, tasks.NewDelayProcessor(backgroundUsecase))

		if err := queue.Run(mux); err != nil {
			if !errors.Is(err, asynq.ErrServerClosed) {
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
}

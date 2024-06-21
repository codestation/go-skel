// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.megpoid.dev/go-skel/app/tasks"
	"go.megpoid.dev/go-skel/app/usecase"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/pkg/cfg"
	"go.megpoid.dev/go-skel/pkg/logger"
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
		logger.InitLogger()
		printVersion()

		generalSettings := config.GeneralSettings{}
		if err := cfg.ReadConfig(&generalSettings); err != nil {
			return fmt.Errorf("failed to read config: %w", err)
		}

		queue := asynq.NewServer(
			asynq.RedisClientOpt{Addr: generalSettings.RedisAddr},
			asynq.Config{Concurrency: generalSettings.Workers},
		)

		backgroundUsecase := usecase.NewDelay()

		mux := asynq.NewServeMux()
		mux.Handle(tasks.TypeDelay, tasks.NewDelayProcessor(backgroundUsecase))

		if err := queue.Run(mux); err != nil {
			if !errors.Is(err, asynq.ErrServerClosed) {
				return fmt.Errorf("failed to run queue: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(queueCmd)

	databaseFlags := config.LoadDatabaseFlags(queueCmd.Name())
	generalFlags := config.LoadGeneralFlags(queueCmd.Name())

	queueCmd.Flags().AddFlagSet(databaseFlags)
	queueCmd.Flags().AddFlagSet(generalFlags)
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/pkg/task"
)

const (
	TypeDelay = "timer:delay"
)

type DelayPayload struct {
	Delay time.Duration
}

func NewDelayTask(delay time.Duration) (*asynq.Task, error) {
	payload, err := json.Marshal(DelayPayload{Delay: delay})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeDelay, payload), nil
}

type DelayTask struct {
	backgroundJob usecase.DelayJob
}

func (process *DelayTask) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p DelayPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	result, err := process.backgroundJob.Process(ctx, p.Delay)
	if err != nil {
		return fmt.Errorf("failed to run job: %w", err)
	}

	response := task.Response{
		ContentType: "application/json",
		Data:        result,
	}

	encoder := json.NewEncoder(t.ResultWriter())
	if err := encoder.Encode(response); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

func NewDelayProcessor(backgroundJob usecase.DelayJob) *DelayTask {
	return &DelayTask{backgroundJob: backgroundJob}
}

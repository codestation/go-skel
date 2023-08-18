// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// used to validate that the implementation matches the interface
var _ Task = &AsynqTask{}

type AsynqTask struct {
	inspector *asynq.Inspector
	client    *asynq.Client
}

func (u *AsynqTask) Enqueue(ctx context.Context, task *asynq.Task) (string, error) {
	info, err := u.client.EnqueueContext(ctx, task, asynq.Queue(DefaultQueueName), asynq.Retention(24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("failed to enqueue task: %w", err)
	}

	return info.ID, err
}

func (u *AsynqTask) GetTaskInfo(_ context.Context, queue, id string) (*Info, error) {
	info, err := u.inspector.GetTaskInfo(queue, id)
	if err != nil {
		return nil, err
	}

	status := &Info{
		ID: info.ID,
	}

	switch info.State {
	case asynq.TaskStateActive:
		status.State = "running"
	case asynq.TaskStatePending:
		fallthrough
	case asynq.TaskStateAggregating:
		fallthrough
	case asynq.TaskStateScheduled:
		status.State = "pending"
	case asynq.TaskStateRetry:
		status.State = "retry"
	case asynq.TaskStateArchived:
		status.State = "failed"
	case asynq.TaskStateCompleted:
		status.State = "succeeded"
	default:
		status.State = "unknown"
	}

	if info.LastErr != "" {
		status.Error = info.LastErr
	}

	return status, nil
}

func (u *AsynqTask) GetTaskResponse(_ context.Context, queue, id string) (*Response, error) {
	info, err := u.inspector.GetTaskInfo(queue, id)
	if err != nil {
		return nil, err
	}

	if info.State != asynq.TaskStateCompleted {
		return nil, fmt.Errorf("task isn't completed yet")
	}

	response := &Response{}
	if err := json.Unmarshal(info.Result, response); err != nil {
		return nil, err
	}

	return response, nil
}

func NewClient(redis asynq.RedisClientOpt) *AsynqTask {
	return &AsynqTask{
		inspector: asynq.NewInspector(redis),
		client:    asynq.NewClient(redis),
	}
}

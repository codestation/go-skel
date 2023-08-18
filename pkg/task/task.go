// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package task

import (
	"context"

	"github.com/hibiken/asynq"
)

const DefaultQueueName = "default"

type Task interface {
	Enqueue(ctx context.Context, task *asynq.Task) (string, error)
	GetTaskInfo(ctx context.Context, queue, id string) (*Info, error)
	GetTaskResponse(ctx context.Context, queue, id string) (*Response, error)
}

type Info struct {
	ID    string `json:"id"`
	State string `json:"status"`
	Error string `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	ContentType string `json:"content_type"`
	Data        any    `json:"data"`
	Error       *Error `json:"error,omitempty"`
}

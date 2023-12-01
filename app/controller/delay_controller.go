// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/app/tasks"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/apperror"
	"go.megpoid.dev/go-skel/pkg/task"
)

type DelayController struct {
	common
	task task.Task
}

func NewDelay(cfg config.ServerSettings, task task.Task) DelayController {
	return DelayController{
		common: newCommon(cfg),
		task:   task,
	}
}

func (ctrl *DelayController) ProcessBackground(ctx echo.Context) error {
	t := ctrl.printer(ctx)

	var request oapi.DelayRequest

	if err := ctx.Bind(&request); err != nil {
		return apperror.NewAppError(t.Sprintf("Failed to read request"), err)
	}

	delay, err := time.ParseDuration(request.Delay)
	if err != nil {
		return apperror.NewValidationError(t.Sprintf("Failed to parse duration"), err)
	}

	delayTask, err := tasks.NewDelayTask(delay)
	if err != nil {
		return apperror.NewValidationError(t.Sprintf("Failed to create task"), err)
	}

	taskId, err := ctrl.task.Enqueue(ctx.Request().Context(), delayTask)
	if err != nil {
		return apperror.NewValidationError(t.Sprintf("Failed to enqueue task"), err)
	}

	response := oapi.TaskCreationResponse{
		Location: "",
		TaskId:   taskId,
	}

	return ctx.JSON(http.StatusAccepted, response)
}

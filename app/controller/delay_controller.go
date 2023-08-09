// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/app/tasks"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/oapi"
	"megpoid.dev/go/go-skel/pkg/apperror"
)

type DelayController struct {
	common
	task usecase.Task
}

func NewDelay(cfg *config.Config, task usecase.Task) DelayController {
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

	task, err := tasks.NewDelayTask(delay)
	if err != nil {
		return apperror.NewValidationError(t.Sprintf("Failed to create task"), err)
	}

	taskId, err := ctrl.task.Enqueue(ctx.Request().Context(), task)
	if err != nil {
		return apperror.NewValidationError(t.Sprintf("Failed to enqueue task"), err)
	}

	response := oapi.TaskCreationResponse{
		Location: "",
		TaskId:   taskId,
	}

	return ctx.JSON(http.StatusAccepted, response)
}
// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/task"
)

type TaskController struct {
	common
	task task.Task
}

func NewTask(cfg config.ServerSettings, task task.Task) TaskController {
	return TaskController{
		common: newCommon(cfg),
		task:   task,
	}
}

func (ctrl *TaskController) GetTask(ctx echo.Context, queueName string, taskId oapi.TaskId) error {
	info, err := ctrl.task.GetTaskInfo(ctx.Request().Context(), queueName, taskId)
	if err != nil {
		return err
	}

	response := oapi.Task{
		State:  oapi.TaskState(info.State),
		TaskId: info.ID,
	}

	if info.Error != "" {
		response.Error = nil
	}

	return ctx.JSON(http.StatusOK, &response)
}

func (ctrl *TaskController) GetTaskResponse(ctx echo.Context, queueName string, taskId oapi.TaskId) error {
	info, err := ctrl.task.GetTaskResponse(ctx.Request().Context(), queueName, taskId)
	if err != nil {
		return err
	}

	switch info.ContentType {
	case echo.MIMEApplicationJSON:
		return ctx.JSON(http.StatusOK, info.Data)
	default:
		return ctx.Blob(http.StatusOK, info.ContentType, info.Data.([]byte))
	}
}

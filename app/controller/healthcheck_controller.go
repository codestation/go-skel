// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/oapi"
)

type HealthcheckController struct {
	common
	healthcheckUsecase usecase.Healthcheck
}

func NewHealthCheck(cfg *config.Config, healthcheck usecase.Healthcheck) HealthcheckController {
	return HealthcheckController{
		common:             newCommon(cfg),
		healthcheckUsecase: healthcheck,
	}
}

func (ctrl *HealthcheckController) LiveCheck(ctx echo.Context, params oapi.LiveCheckParams) error {
	if params.Verbose != nil && *params.Verbose {
		var check strings.Builder
		check.WriteString("live check passed\n")
		return ctx.String(http.StatusOK, check.String())
	}
	return ctx.String(http.StatusOK, "ok")
}

func (ctrl *HealthcheckController) ReadyCheck(ctx echo.Context, params oapi.ReadyCheckParams) error {
	result := ctrl.healthcheckUsecase.Execute(ctx.Request().Context())
	if params.Verbose != nil && *params.Verbose {
		var check strings.Builder
		if result.Ping != nil {
			check.WriteString(fmt.Sprintf("[+] ping err: %s\n", result.Ping.Error()))
		} else {
			check.WriteString("[+] ping ok\n")
		}
		if result.AllOk() {
			check.WriteString("ready check passed\n")
		} else {
			check.WriteString("ready check failed\n")
		}
		return ctx.String(http.StatusOK, check.String())
	}

	if !result.AllOk() {
		return echo.NewHTTPError(http.StatusInternalServerError, "error")
	}

	return ctx.String(http.StatusOK, "ok")
}

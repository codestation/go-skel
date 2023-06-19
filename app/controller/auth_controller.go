// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/oapi"
)

type AuthController struct {
	common
	auth usecase.Auth
}

func NewAuth(cfg *config.Config, Auth usecase.Auth) AuthController {
	return AuthController{
		common: newCommon(cfg),
		auth:   Auth,
	}
}

func (ctrl *AuthController) Login(ctx echo.Context) error {
	t := ctrl.printer(ctx)

	request := &oapi.AuthRequest{}

	if err := ctx.Bind(request); err != nil {
		return usecase.NewAppError(t.Sprintf("Invalid login"), err)
	}
	if err := ctx.Validate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	result, err := ctrl.auth.Login(ctx.Request().Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &oapi.Token{Token: result})
}

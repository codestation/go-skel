// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/app/usecase"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/apperror"
	"go.megpoid.dev/go-skel/pkg/middleware"
)

type AuthController struct {
	common
	auth usecase.Auth
	oidc *middleware.Auth
}

func NewAuth(cfg config.ServerSettings, auth usecase.Auth, oidc *middleware.Auth) AuthController {
	return AuthController{
		common: newCommon(cfg),
		auth:   auth,
		oidc:   oidc,
	}
}

func (ctrl *AuthController) Login(ctx echo.Context) error {
	t := ctrl.printer(ctx)

	request := oapi.AuthRequest{}

	if err := ctx.Bind(&request); err != nil {
		return apperror.NewValidationError(t.Sprintf("Invalid login request"), err)
	}
	if err := ctx.Validate(&request); err != nil {
		return err
	}

	result, err := ctrl.auth.Login(ctx.Request().Context(), request.Username, request.Password)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &oapi.Token{Token: result})
}

func (ctrl *AuthController) OAuthLogin(ctx echo.Context) error {
	if err := ctrl.oidc.RedirectHandler(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "oAuth redirect")
}

func (ctrl *AuthController) OAuthRefresh(ctx echo.Context) error {
	if err := ctrl.oidc.RefreshHandler(ctx); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, "token refreshed")
}

func (ctrl *AuthController) OAuthCallback(ctx echo.Context) error {
	if err := ctrl.oidc.CallbackHandler(ctx); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(http.StatusFound, "/")
}

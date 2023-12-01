// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/app/controller/filter"
	"go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/app/usecase"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/apperror"
)

type ProfileController struct {
	common
	profileUsecase usecase.Profile
}

func NewProfile(cfg config.ServerSettings, profile usecase.Profile) ProfileController {
	return ProfileController{
		common:         newCommon(cfg),
		profileUsecase: profile,
	}
}

func (ctrl *ProfileController) ListProfiles(ctx echo.Context, params oapi.ListProfilesParams) error {
	query, err := filter.NewFilterFromParams(filter.Params(params))
	if err != nil {
		return err
	}

	result, err := ctrl.profileUsecase.ListProfiles(ctx.Request().Context(), query)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (ctrl *ProfileController) SaveProfile(ctx echo.Context) error {
	t := ctrl.printer(ctx)

	var request oapi.ProfileRequest
	if err := ctx.Bind(&request); err != nil {
		return apperror.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := ctx.Validate(&request); err != nil {
		return apperror.NewValidationError(t.Sprintf("The request did not pass validation"), err)
	}

	result, err := ctrl.profileUsecase.SaveProfile(ctx.Request().Context(), (*model.ProfileRequest)(&request))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (ctrl *ProfileController) RemoveProfile(ctx echo.Context, id oapi.ProfileId) error {
	err := ctrl.profileUsecase.RemoveProfile(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (ctrl *ProfileController) GetProfile(ctx echo.Context, id oapi.ProfileId) error {
	result, err := ctrl.profileUsecase.GetProfile(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (ctrl *ProfileController) UpdateProfile(ctx echo.Context, id oapi.ProfileId) error {
	var request oapi.ProfileRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	if err := ctx.Validate(&request); err != nil {
		return err
	}

	result, err := ctrl.profileUsecase.UpdateProfile(ctx.Request().Context(), id, (*model.ProfileRequest)(&request))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

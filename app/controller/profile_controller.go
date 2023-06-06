// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/app/controller/filter"
	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/oapi"
)

type ProfileController struct {
	common
	profile usecase.Profile
}

func NewProfileCtrl(cfg *config.Config, profile usecase.Profile) ProfileController {
	return ProfileController{
		common:  newCommon(cfg),
		profile: profile,
	}
}

func (a *ProfileController) ListProfiles(ctx echo.Context, params oapi.ListProfilesParams) error {
	query, err := filter.NewFilterFromParams(filter.Params(params))
	if err != nil {
		return err
	}

	result, err := a.profile.ListProfiles(ctx.Request().Context(), query)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (a *ProfileController) SaveProfile(ctx echo.Context) error {
	t := a.printer(ctx)

	var request oapi.ProfileRequest
	if err := ctx.Bind(&request); err != nil {
		return usecase.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := ctx.Validate(&request); err != nil {
		return usecase.NewAppError(t.Sprintf("The request did not pass validation"), err)
	}

	result, err := a.profile.SaveProfile(ctx.Request().Context(), (*model.ProfileRequest)(&request))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, result)
}

func (a *ProfileController) RemoveProfile(ctx echo.Context, id oapi.ProfileId) error {
	err := a.profile.RemoveProfile(ctx.Request().Context(), model.ID(id))
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (a *ProfileController) GetProfile(ctx echo.Context, id oapi.ProfileId) error {
	result, err := a.profile.GetProfile(ctx.Request().Context(), model.ID(id))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (a *ProfileController) UpdateProfile(ctx echo.Context, id oapi.ProfileId) error {
	var request oapi.ProfileRequest
	if err := ctx.Bind(&request); err != nil {
		return err
	}
	if err := ctx.Validate(&request); err != nil {
		return err
	}

	result, err := a.profile.UpdateProfile(ctx.Request().Context(), model.ID(id), (*model.ProfileRequest)(&request))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.dev/go/go-skel/app"
	"megpoid.dev/go/go-skel/app/i18n"
	"megpoid.dev/go/go-skel/model"
)

func (api *API) InitProfile() {
	protected := api.root.Group("")
	protected.GET("/profiles", api.ListProfiles)
	protected.GET("/profiles/:id", api.GetProfile)
	protected.POST("/profiles", api.SaveProfile)
	protected.PATCH("/profiles/:id", api.UpdateProfile)
	protected.DELETE("/profiles/:id", api.RemoveProfile)
}

func (api *API) GetProfile(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var profileId int
	err := echo.PathParamsBinder(c).MustInt("id", &profileId).BindError()
	if err != nil {
		return app.NewAppError(t.Sprintf("Invalid profile ID"), err)
	}

	result, err := api.app.GetProfile(c.Request().Context(), model.ID(profileId))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (api *API) ListProfiles(c echo.Context) error {
	query, err := NewFilter(c)
	if err != nil {
		return err
	}

	result, err := api.app.ListProfiles(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (api *API) SaveProfile(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var request model.ProfileRequest
	if err := c.Bind(&request); err != nil {
		return app.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := c.Validate(&request); err != nil {
		return app.NewAppError(t.Sprintf("The request did not pass validation"), err)
	}

	result, err := api.app.SaveProfile(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
func (api *API) UpdateProfile(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var request model.ProfileRequest
	if err := c.Bind(&request); err != nil {
		return app.NewAppError(t.Sprintf("Failed to read request"), err)
	}
	if err := c.Validate(&request); err != nil {
		return app.NewAppError(t.Sprintf("The request did not pass validation"), err)
	}

	result, err := api.app.UpdateProfile(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
func (api *API) RemoveProfile(c echo.Context) error {
	t := message.NewPrinter(i18n.GetLanguageTags(c))

	var profileId int
	err := echo.PathParamsBinder(c).MustInt("id", &profileId).BindError()
	if err != nil {
		return app.NewAppError(t.Sprintf("Invalid profile ID"), err)
	}

	err = api.app.RemoveProfile(c.Request().Context(), model.ID(profileId))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

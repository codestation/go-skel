// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.xyz/go/go-skel/app"
	"megpoid.xyz/go/go-skel/app/i18n"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/response"
	"net/http"
)

func (api *API) InitProfile() {
	protected := api.root.Group("")
	protected.GET("/profiles", api.ListProfiles)
	protected.GET("/profiles/:id", api.GetProfile)
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

	results, cur, err := api.app.ListProfiles(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.NewListResponse(results, cur))
}

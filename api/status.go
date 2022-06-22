// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (api *API) InitStatus() {
	api.root.GET("/status/livez", api.LiveCheck)
	api.root.GET("/status/readyz", api.ReadyCheck)
}

func (api *API) LiveCheck(c echo.Context) error {
	if c.QueryParams().Has("verbose") {
		var check strings.Builder
		check.WriteString("livez check passed\n")
		return c.String(http.StatusOK, check.String())
	}
	return c.String(http.StatusOK, "ok")
}

func (api *API) ReadyCheck(c echo.Context) error {
	result := api.app.HealthCheck(c.Request().Context())
	if c.QueryParams().Has("verbose") {
		var check strings.Builder
		if result.Ping != nil {
			check.WriteString(fmt.Sprintf("[+] ping err: %s\n", result.Ping.Error()))
		} else {
			check.WriteString("[+] ping ok\n")
		}
		if result.AllOk() {
			check.WriteString("readyz check passed\n")
		} else {
			check.WriteString("readyz check failed\n")
		}
		return c.String(http.StatusOK, check.String())
	} else {
		if !result.AllOk() {
			return echo.NewHTTPError(http.StatusInternalServerError, "error")
		}
		return c.String(http.StatusOK, "ok")
	}
}

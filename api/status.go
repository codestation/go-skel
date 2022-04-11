// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

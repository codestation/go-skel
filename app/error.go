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

package app

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"megpoid.xyz/go/go-skel/i18n"
	"megpoid.xyz/go/go-skel/model"
)

const (
	TranslateKey = "translationFunc"
)

func ErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appErr *model.AppError
		var httpErr *echo.HTTPError
		T := c.Get(TranslateKey).(i18n.TranslateFunc)

		switch {
		case errors.As(err, &appErr):
			appErr.Translate(T)
			// remove internal state from non-debug error response
			if !e.Debug {
				appErr.DetailedError = ""
				appErr.Location = ""
			}
		case errors.As(err, &httpErr):
			appErr = model.NewAppError("", "app_error", nil, "", httpErr.Code)
			appErr.Translate(T)
			if e.Debug {
				if httpErr.Internal != nil {
					appErr.DetailedError = httpErr.Internal.Error()
				} else {
					appErr.DetailedError = http.StatusText(httpErr.Code)
				}
				appErr.Location = "app"
			}
		default:
			appErr = model.NewAppError("app", "app_error", nil, err.Error(), http.StatusInternalServerError)
			appErr.Translate(T)
		}

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(appErr.StatusCode)
		} else {
			err = c.JSON(appErr.StatusCode, appErr)
		}
		if err != nil {
			e.Logger.Error(err)
		}
	}
}

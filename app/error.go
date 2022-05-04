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
	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.xyz/go/go-skel/app/i18n"
	"megpoid.xyz/go/go-skel/model"
	"net/http"
)

func ErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appErr *model.AppError
		var httpErr *echo.HTTPError
		printer := message.NewPrinter(i18n.GetLanguageTags(c))

		switch {
		case errors.As(err, &appErr):
			// remove internal state from non-debug error response
			if !e.Debug {
				appErr.DetailedError = ""
				appErr.Location = ""
			}
		case errors.As(err, &httpErr):
			appErr = model.NewAppError("app", printer.Sprintf("An error occurred"), "", httpErr.Code)
			if e.Debug {
				if httpErr.Internal != nil {
					appErr.DetailedError = httpErr.Internal.Error()
				} else {
					appErr.DetailedError = http.StatusText(httpErr.Code)
				}
			} else {
				appErr.Location = ""
			}
		default:
			appErr = model.NewAppError("app", printer.Sprintf("An error occurred"), err.Error(), http.StatusInternalServerError)
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

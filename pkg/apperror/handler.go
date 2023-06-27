// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package apperror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.dev/go/go-skel/pkg/i18n"
)

func ErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var appErr *Error
		printer := message.NewPrinter(i18n.GetLanguageTags(c))

		if !errors.As(err, &appErr) {
			appErr = NewAppError(printer.Sprintf("An error occurred"), err)
			appErr.Where = "ErrorHandler"
		}

		if e.Debug {
		} else {
			appErr.DetailedError = ""
			appErr.Where = ""
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

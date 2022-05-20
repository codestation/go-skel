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
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"golang.org/x/text/message"
	"megpoid.xyz/go/go-skel/app/i18n"
	"megpoid.xyz/go/go-skel/store"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Error struct {
	Message       string            `json:"message"`
	Where         string            `json:"location,omitempty"`
	DetailedError string            `json:"detailed_error,omitempty"`
	StatusCode    int               `json:"status_code"`
	Validation    []ValidationError `json:"validation,omitempty"`
}

func (e *Error) Error() string {
	return e.Where + ": " + e.Message + ", " + e.DetailedError
}

func NewAppError(message string, err error) *Error {
	appErr := &Error{
		Message: message,
	}

	pc := make([]uintptr, 1)
	n := runtime.Callers(2, pc)

	if n > 0 {
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		funcParts := strings.Split(frame.Function, ".")
		appErr.Where = funcParts[len(funcParts)-1]
	}

	if err != nil {
		var httpErr *echo.HTTPError
		var validateErr validator.ValidationErrors
		var bindingErr *echo.BindingError

		appErr.DetailedError = err.Error()
		switch {
		case errors.Is(err, store.ErrNotFound):
			appErr.StatusCode = http.StatusNotFound
		case errors.As(err, &httpErr):
			appErr.StatusCode = httpErr.Code
		case errors.As(err, &bindingErr):
			appErr.StatusCode = bindingErr.Code
			appErr.DetailedError = bindingErr.Internal.Error()
		case errors.As(err, &validateErr):
			for _, v := range validateErr {
				appErr.Validation = append(appErr.Validation, ValidationError{
					Field:   v.Field(),
					Message: v.ActualTag(),
				})
			}
			appErr.StatusCode = http.StatusBadRequest
		default:
			appErr.StatusCode = http.StatusInternalServerError
		}
	} else {
		appErr.StatusCode = http.StatusInternalServerError
	}

	return appErr
}

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

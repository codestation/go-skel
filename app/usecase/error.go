// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"errors"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/pkg/repo"
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
		case errors.Is(err, repo.ErrNotFound):
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

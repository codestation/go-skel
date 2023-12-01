// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package apperror

import (
	"errors"
	"net/http"
	"runtime"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/pkg/repo"
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

func createAppError(message string, err error, skip int) *Error {
	appErr := &Error{
		Message: message,
	}

	pc := make([]uintptr, 1)
	n := runtime.Callers(skip, pc)

	if n > 0 {
		frames := runtime.CallersFrames(pc[:n])
		frame, _ := frames.Next()
		appErr.Where = frame.Function
	}

	if err != nil {
		var httpErr *echo.HTTPError
		var validateErr validator.ValidationErrors
		var bindingErr *echo.BindingError
		var authnErr *authnError
		var authzErr *authzError
		var validationErr *validationError

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
		case errors.As(err, &authnErr):
			appErr.StatusCode = http.StatusUnauthorized
			appErr.DetailedError = authnErr.Error()
		case errors.As(err, &validationErr):
			appErr.StatusCode = http.StatusBadRequest
			appErr.DetailedError = validationErr.Error()
		case errors.As(err, &authzErr):
			appErr.StatusCode = http.StatusForbidden
			appErr.DetailedError = authzErr.Error()
		default:
			appErr.StatusCode = http.StatusInternalServerError
		}
	} else {
		appErr.StatusCode = http.StatusInternalServerError
	}

	return appErr
}

func NewAppError(message string, err error) *Error {
	return createAppError(message, err, 3)
}

type authnError struct {
	err error
}

func (e *authnError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

type validationError struct {
	err error
}

func (e *validationError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func NewAuthnError(message string, err error) *Error {
	return createAppError(message, &authnError{err: err}, 4)
}

func NewValidationError(message string, err error) *Error {
	return createAppError(message, &validationError{err: err}, 4)
}

type authzError struct {
	err error
}

func (e *authzError) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

func NewAuthzError(message string, err error) *Error {
	return createAppError(message, &authzError{err: err}, 4)
}

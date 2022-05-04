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

package model

import (
	"megpoid.xyz/go/go-skel/i18n"
)

type AppError struct {
	ID            string `json:"id"`
	Message       string `json:"message"`
	Location      string `json:"location,omitempty"`
	DetailedError string `json:"detailed_error,omitempty"`
	StatusCode    int    `json:"status_code"`
	params        map[string]any
}

func (e *AppError) Error() string {
	return e.Location + ": " + e.Message + ", " + e.DetailedError
}

func (e *AppError) Translate(T i18n.TranslateFunc) {
	if T == nil {
		e.Message = e.ID
	} else {
		e.Message = T(e.ID, e.params)
	}
}

func NewAppError(location, id string, params map[string]any, details string, status int) *AppError {
	appErr := &AppError{
		ID:            id,
		Location:      location,
		DetailedError: details,
		StatusCode:    status,
		params:        params,
	}

	appErr.Translate(i18n.T)
	return appErr
}

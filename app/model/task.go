// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

type Task struct {
	ID    string `json:"id"`
	State string `json:"status"`
	Error string `json:"error"`
}

type TaskError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type TaskResponse struct {
	ContentType string     `json:"content_type"`
	Data        any        `json:"data"`
	Error       *TaskError `json:"error,omitempty"`
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"errors"
	"fmt"
)

type RepoError struct {
	Err      error
	internal error
}

var ErrBackend = errors.New("repo: backend error")
var ErrNotFound = errors.New("repo: model not found")
var ErrDuplicated = errors.New("repo: duplicated model")

func NewRepoError(err, internal error) error {
	return &RepoError{internal: internal, Err: err}
}

func (r *RepoError) Error() string {
	if r.internal != nil {
		return fmt.Sprintf("%s: %s", r.Err, r.internal)
	}
	return r.Err.Error()
}

func (r *RepoError) Unwrap() error {
	return r.Err
}

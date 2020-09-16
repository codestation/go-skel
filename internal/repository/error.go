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
var ErrNotFound = errors.New("repo: entity not found")
var ErrDuplicated = errors.New("repo: duplicated entity")

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

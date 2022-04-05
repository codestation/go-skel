/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store

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

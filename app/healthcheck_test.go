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
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/store"
)

type fakeStore struct {
	HealthCheckError error
}

func (s fakeStore) HealthCheck() store.HealthCheckStore {
	return &fakeHealthCheckStore{
		Error: s.HealthCheckError,
	}
}

func (s fakeStore) WithTransaction(_ context.Context, _ func(s store.Store) error) error {
	panic("implement me")
}

type fakeHealthCheckStore struct {
	Error error
}

func (f fakeHealthCheckStore) HealthCheck(_ context.Context) error {
	return f.Error
}

func TestApp_HealthCheck(t *testing.T) {
	ss := &fakeStore{}
	srv := &Server{Store: ss}
	app := New(srv)
	result := app.HealthCheck(context.Background())
	assert.NotNil(t, result)
	assert.True(t, result.AllOk())
}

func TestApp_HealthCheckError(t *testing.T) {
	ss := &fakeStore{HealthCheckError: errors.New("an error")}
	srv := &Server{Store: ss}
	app := New(srv)
	result := app.HealthCheck(context.Background())
	assert.NotNil(t, result)
	assert.False(t, result.AllOk())
}

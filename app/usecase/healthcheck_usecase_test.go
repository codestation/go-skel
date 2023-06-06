// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"megpoid.dev/go/go-skel/app/repository"
)

func TestApp_Healthcheck(t *testing.T) {
	repo := repository.NewMockHealthcheckRepo(t)
	repo.EXPECT().Execute(mock.Anything).Return(nil)

	u := NewHealthcheck(repo)
	result := u.Execute(context.Background())
	assert.NotNil(t, result)
	assert.True(t, result.AllOk())
}

func TestApp_HealthcheckError(t *testing.T) {
	repo := repository.NewMockHealthcheckRepo(t)
	repo.EXPECT().Execute(mock.Anything).Return(errors.New("an error"))

	app := NewHealthcheck(repo)
	result := app.Execute(context.Background())
	assert.NotNil(t, result)
	assert.False(t, result.AllOk())
}

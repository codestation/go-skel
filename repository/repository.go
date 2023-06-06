// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/repo"
)

// HealthcheckRepo handles all healthCheck related operations on the repository
//
//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name HealthcheckRepo
type HealthcheckRepo interface {
	Execute(ctx context.Context) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name ProfileRepo
type ProfileRepo interface {
	repo.GenericStore[*model.Profile]
	GetByEmail(ctx context.Context, email string) (*model.Profile, error)
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/pkg/repo"
)

// HealthcheckRepo handles all healthCheck related operations on the repository
type HealthcheckRepo interface {
	Execute(ctx context.Context) error
}

type ProfileRepo interface {
	repo.GenericStore[*model.Profile]
	GetByEmail(ctx context.Context, email string) (*model.Profile, error)
}

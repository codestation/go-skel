// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"log/slog"

	"megpoid.dev/go/go-skel/app/repository"
)

// used to validate that the implementation matches the interface
var _ Healthcheck = &HealthcheckInteractor{}

type HealthcheckInteractor struct {
	common
	healthcheckRepo repository.HealthcheckRepo
}

func (u *HealthcheckInteractor) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Executing healthcheck")
	return u.healthcheckRepo.Execute(ctx)
}

func NewHealthcheck(repo repository.HealthcheckRepo) *HealthcheckInteractor {
	return &HealthcheckInteractor{
		common:          newCommon(),
		healthcheckRepo: repo,
	}
}

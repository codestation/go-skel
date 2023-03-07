// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/repository"
)

// used to validate that the implementation matches the interface
var _ Healthcheck = &HealthcheckInteractor{}

type HealthcheckInteractor struct {
	common
	healthcheckRepo repository.HealthcheckRepo
}

func (u *HealthcheckInteractor) Execute(ctx context.Context) *model.HealthcheckResult {
	err := u.healthcheckRepo.Execute(ctx)
	return &model.HealthcheckResult{
		Ping: err,
	}
}

func NewHealthcheck(repo repository.HealthcheckRepo) *HealthcheckInteractor {
	return &HealthcheckInteractor{
		common:          newCommon(),
		healthcheckRepo: repo,
	}
}

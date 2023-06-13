// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/request"
	"megpoid.dev/go/go-skel/pkg/response"
)

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name Profile
type Profile interface {
	GetProfile(ctx context.Context, id int64) (*model.Profile, error)
	ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[*model.Profile], error)
	SaveProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error)
	UpdateProfile(ctx context.Context, id int64, req *model.ProfileRequest) (*model.Profile, error)
	RemoveProfile(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name Healthcheck
type Healthcheck interface {
	Execute(ctx context.Context) *model.HealthcheckResult
}

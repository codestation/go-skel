// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"time"

	"go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/pkg/request"
	"go.megpoid.dev/go-skel/pkg/response"
)

type Auth interface {
	Login(ctx context.Context, username, password string) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name Profile
type Profile interface {
	GetProfile(ctx context.Context, id int64) (*model.Profile, error)
	ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[*model.Profile], error)
	SaveProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error)
	UpdateProfile(ctx context.Context, id int64, req *model.ProfileRequest) (*model.Profile, error)
	RemoveProfile(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name Healthcheck
type Healthcheck interface {
	Execute(ctx context.Context) error
}

type DelayJob interface {
	Process(ctx context.Context, delay time.Duration) (*Timers, error)
}

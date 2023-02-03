// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/model"
	"megpoid.dev/go/go-skel/model/request"
	"megpoid.dev/go/go-skel/model/response"
)

type IApp interface {
	GetProfile(ctx context.Context, id model.ID) (*model.Profile, error)
	ListProfiles(ctx context.Context, query *request.QueryParams) (*response.ListResponse[*model.Profile], error)
	SaveProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error)
	UpdateProfile(ctx context.Context, req *model.ProfileRequest) (*model.Profile, error)
	RemoveProfile(ctx context.Context, id model.ID) error
	HealthCheck(ctx context.Context) *model.HealthCheckResult
	Srv() *Server
	Config() *config.Config
}

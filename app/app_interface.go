// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"megpoid.xyz/go/go-skel/config"
	"megpoid.xyz/go/go-skel/model/request"
	"megpoid.xyz/go/go-skel/store/paginator/cursor"

	"megpoid.xyz/go/go-skel/model"
)

type IApp interface {
	GetProfile(ctx context.Context, id model.ID) (*model.Profile, error)
	ListProfiles(ctx context.Context, query *request.QueryParams) ([]*model.Profile, *cursor.Cursor, error)

	HealthCheck(ctx context.Context) *model.HealthCheckResult
	Srv() *Server
	Config() *config.Config
}

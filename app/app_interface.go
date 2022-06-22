// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"megpoid.xyz/go/go-skel/config"

	"megpoid.xyz/go/go-skel/model"
)

type IApp interface {
	HealthCheck(ctx context.Context) *model.HealthCheckResult
	Srv() *Server
	Config() *config.Config
}

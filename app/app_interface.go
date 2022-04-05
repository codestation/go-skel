package app

import (
	"context"

	"megpoid.xyz/go/go-skel/model"
)

type IApp interface {
	HealthCheck(ctx context.Context) *model.HealthCheckResult
	Srv() *Server
	Config() *model.Config
}

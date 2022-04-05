package app

import (
	"context"

	"megpoid.xyz/go/go-skel/model"
)

func (a *App) HealthCheck(ctx context.Context) *model.HealthCheckResult {
	err := a.Srv().Store.HealthCheck().HealthCheck(ctx)
	return &model.HealthCheckResult{
		Ping: err,
	}
}

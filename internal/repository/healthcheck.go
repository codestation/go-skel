package repository

import (
	"context"

	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

type pingRepo struct {
	db connection.SQLConnection
}

// NewHealthCheck creates a new repository with methods to check its health
func NewHealthCheck(db connection.SQLConnection) HealthCheck {
	return &pingRepo{db}
}

// HealthCheck returns an error if the database doesn't respond
func (r *pingRepo) HealthCheck(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return NewRepoError(ErrBackend, err)
	}
	return nil
}

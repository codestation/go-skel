//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

package repository

import (
	"context"

	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

// HealthCheck handles all health-related operations on the repository
type HealthCheck interface {
	HealthCheck(ctx context.Context) error
}

// Transaction allow to call multiple repository methods and make their execution atomic
type Transaction interface {
	WithTransaction(ctx context.Context, txFunc func(repo *Repository) error) error
}

// Repository holds the rest of the repository interfaces
type Repository struct {
	HealthCheck
	Transaction
}

// NewRepository created a new repository for the connection
func NewRepository(db connection.SQLConnection) *Repository {
	return &Repository{
		HealthCheck: NewHealthCheck(db),
		Transaction: NewTransaction(db),
	}
}

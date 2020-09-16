package repository

import (
	"context"

	"megpoid.xyz/go/go-skel/pkg/sql/connection"
)

type transaction struct {
	db connection.SQLConnection
}

// NewTransaction creates a new transaction handler for the connection
func NewTransaction(db connection.SQLConnection) Transaction {
	return &transaction{db: db}
}

// WithTransaction run a transaction for each repository
func (r *transaction) WithTransaction(ctx context.Context, txFunc func(repo *Repository) error) error {
	tx, err := r.db.TxBegin(ctx)
	if err != nil {
		return err
	}

	return tx.TxEnd(func() error {
		return txFunc(NewRepository(tx))
	})
}

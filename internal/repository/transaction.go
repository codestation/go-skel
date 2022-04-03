/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

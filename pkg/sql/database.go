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

package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"megpoid.xyz/go/go-skel/pkg/sql/helper"
)

type Option func(db *sqlx.DB) error

func MaxOpenConns(n int) Option {
	return func(db *sqlx.DB) error {
		db.SetMaxOpenConns(n)
		return nil
	}
}

// NewDatabase opens a database connection indicated by the dsn connection string
func NewDatabase(ctx context.Context, driverName string, dsn string, opts ...Option) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	db.MapperFunc(helper.ToSnakeCase)

	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

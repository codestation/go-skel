// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"megpoid.xyz/go/go-skel/model"
)

const (
	pingMaxAttempts             = 6
	pingTimeoutSecs             = 10
	postgresUniqueViolationCode = "23505"
)

// compile time validator for the interfaces
var (
	_ SqlExecutor = pgxWrapper{}
	_ SqlExecutor = pgxTxWrapper{}

	ErrNoRows = pgx.ErrNoRows
)

type pgxWrapper struct {
	*pgxpool.Pool
}

func (p pgxWrapper) Begin(ctx context.Context, f func(db SqlExecutor) error) error {
	return p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(newPgxTxWrapper(tx))
	})
}

func (p pgxWrapper) Exec(ctx context.Context, query string, arguments ...interface{}) (sql.Result, error) {
	tag, err := p.Pool.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

func (p pgxWrapper) Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, p, dst, query, args...)
}

func (p pgxWrapper) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, p, dest, query, args...)
}

type pgxTxWrapper struct {
	pgx.Tx
}

func (p pgxTxWrapper) Begin(ctx context.Context, f func(db SqlExecutor) error) error {
	return p.Tx.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(pgxTxWrapper{tx})
	})
}

func (p pgxTxWrapper) Exec(ctx context.Context, query string, arguments ...interface{}) (sql.Result, error) {
	tag, err := p.Tx.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

func (p pgxTxWrapper) Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	return pgxscan.Get(ctx, p, dst, query, args...)
}

func (p pgxTxWrapper) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return pgxscan.Select(ctx, p, dest, query, args...)
}

type pgxWrapperResult struct {
	pgconn.CommandTag
}

func (r pgxWrapperResult) LastInsertId() (int64, error) {
	return 0, errors.New("not implemented")
}

func (r pgxWrapperResult) RowsAffected() (int64, error) {
	return r.CommandTag.RowsAffected(), nil
}

func newPgxWrapper(pool *pgxpool.Pool) *pgxWrapper {
	return &pgxWrapper{pool}
}

func newPgxTxWrapper(tx pgx.Tx) *pgxTxWrapper {
	return &pgxTxWrapper{tx}
}

func NewConnection(settings model.SqlSettings) SqlDb {
	config, err := pgxpool.ParseConfig(settings.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to configure database, aborting: %s", err.Error())
	}

	config.MaxConnLifetime = settings.ConnMaxLifetime
	config.MaxConnIdleTime = settings.ConnMaxIdleTime
	config.MaxConns = int32(settings.MaxOpenConns)
	config.MinConns = int32(settings.MaxIdleConns)

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to open database, aborting: %s", err.Error())
	}

	// total waiting time: 1 minute
	for i := 0; i < pingMaxAttempts; i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutSecs*time.Second)
			defer cancel()

			return db.Ping(ctx)
		}()

		if err == nil {
			break
		}

		if i < pingMaxAttempts {
			log.Printf("Failed to ping database: %s, retrying in %d seconds", err.Error(), pingTimeoutSecs)
			time.Sleep(pingTimeoutSecs * time.Second)
		} else {
			log.Fatal("Failed to ping database, aborting")
		}
	}

	return newPgxWrapper(db)
}

func NewQueryBuilder() goqu.DialectWrapper {
	return goqu.Dialect("postgres")
}

func IsUniqueError(err error, opts ...Option) bool {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		if pgErr.Code == postgresUniqueViolationCode {
			for _, opt := range opts {
				if !opt(pgErr) {
					return false
				}
			}
			return true
		}
	}

	return false
}

type Option func(err error) bool

func WithConstraintName(name string) Option {
	return func(err error) bool {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):
			if pgErr.ConstraintName == name {
				return true
			}
		}

		return false
	}
}

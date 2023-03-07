// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"megpoid.dev/go/go-skel/config"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	pingMaxAttempts             = 6
	pingTimeoutSecs             = 10
	postgresUniqueViolationCode = "23505"
)

// compile time validator for the interfaces
var (
	_ SqlExecutor = PgxWrapper{}
	_ SqlExecutor = PgxTxWrapper{}

	ErrNoRows = pgx.ErrNoRows
)

type PgxWrapper struct {
	pool *pgxpool.Pool
}

func (p PgxWrapper) BeginFunc(ctx context.Context, f func(conn SqlExecutor) error) error {
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		return f(newPgxTxWrapper(tx))
	})
}

func (p PgxWrapper) Begin(ctx context.Context) (*PgxTxWrapper, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPgxTxWrapper(tx), nil
}

func (p PgxWrapper) Exec(ctx context.Context, query string, arguments ...any) (sql.Result, error) {
	tag, err := p.pool.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

func (p PgxWrapper) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.pool, dst, query, args...)
}

func (p PgxWrapper) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.pool, dest, query, args...)
}

func (p PgxWrapper) Close() {
	p.pool.Close()
}

func (p PgxWrapper) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

type PgxTxWrapper struct {
	tx pgx.Tx
}

func (p PgxTxWrapper) BeginFunc(ctx context.Context, f func(conn SqlExecutor) error) error {
	return pgx.BeginFunc(ctx, p.tx, func(tx pgx.Tx) error {
		return f(PgxTxWrapper{tx})
	})
}

func (p PgxTxWrapper) Begin(ctx context.Context) (*PgxTxWrapper, error) {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPgxTxWrapper(tx), nil
}

func (p PgxTxWrapper) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

func (p PgxTxWrapper) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}

func (p PgxTxWrapper) Exec(ctx context.Context, query string, arguments ...any) (sql.Result, error) {
	tag, err := p.tx.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

func (p PgxTxWrapper) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.tx, dst, query, args...)
}

func (p PgxTxWrapper) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.tx, dest, query, args...)
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

func NewPgxWrapper(pool *pgxpool.Pool) *PgxWrapper {
	return &PgxWrapper{pool}
}

func newPgxTxWrapper(tx pgx.Tx) *PgxTxWrapper {
	return &PgxTxWrapper{tx}
}

func NewConnection(settings config.SqlSettings) (*pgxpool.Pool, error) {
	parseConfig, err := pgxpool.ParseConfig(settings.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to configure database, aborting: %w", err)
	}

	parseConfig.MaxConnLifetime = settings.ConnMaxLifetime
	parseConfig.MaxConnIdleTime = settings.ConnMaxIdleTime
	parseConfig.MaxConns = int32(settings.MaxOpenConns)
	parseConfig.MinConns = int32(settings.MaxIdleConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database, aborting: %w", err)
	}

	// total waiting time: 1 minute
	for i := 0; i < pingMaxAttempts; i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutSecs*time.Second)
			defer cancel()

			return pool.Ping(ctx)
		}()

		if err == nil {
			break
		}

		if i < pingMaxAttempts {
			log.Printf("Failed to ping database: %s, retrying in %d seconds", err.Error(), pingTimeoutSecs)
			time.Sleep(pingTimeoutSecs * time.Second)
		} else {
			return nil, errors.New("failed to ping database, aborting")
		}
	}

	goqu.SetColumnRenameFunction(dbscan.SnakeCaseMapper)

	return pool, nil
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

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name Querier
//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --srcpkg github.com/jackc/pgx/v5 --inpackage=false --output . --outpkg=sqlrepo --filename=rows_mock.go --name Rows --structname MockRows
type Querier interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

func GetStruct[T any](ctx context.Context, p Querier, query string, args ...any) (T, error) {
	rows, err := p.Query(ctx, query, args...)
	if err != nil {
		var t T
		return t, err
	}

	return pgx.CollectOneRow[T](rows, pgx.RowToStructByName[T])
}

func SelectStruct[T any](ctx context.Context, p Querier, query string, args ...any) ([]T, error) {
	rows, err := p.Query(ctx, query, args...)
	if err != nil {
		var t []T
		return t, err
	}

	return pgx.CollectRows[T](rows, pgx.RowToStructByName[T])
}

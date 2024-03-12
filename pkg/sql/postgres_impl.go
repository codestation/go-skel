// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	// import postgres dialect for goqu library
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pingMaxAttempts             = 5
	pingTimeoutSecs             = 10
	postgresUniqueViolationCode = "23505"
)

// compile time validator for the interfaces
var (
	_ Executor = &PgxPool{}
	_ Executor = &PgxTx{}
)

// PgxPool is a PostgreSQL wrapper that implements the Executor interface.
type PgxPool struct {
	*pgxpool.Pool
}

func NewPgxPool(pool *pgxpool.Pool) *PgxPool {
	return &PgxPool{pool}
}

func (p *PgxPool) Begin(ctx context.Context) (*PgxTx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return NewPgxTx(tx), nil
}

func (p *PgxPool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (*PgxTx, error) {
	tx, err := p.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, err
	}

	return NewPgxTx(tx), nil
}

// BeginFunc starts a transaction and executes the given function within that transaction.
func (p *PgxPool) BeginFunc(ctx context.Context, f func(conn Tx) error) error {
	return pgx.BeginFunc(ctx, p.Pool, func(tx pgx.Tx) error {
		return f(NewPgxTx(tx))
	})
}

// BeginTxFunc starts a transaction and executes the given function within that transaction.
func (p *PgxPool) BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(conn Tx) error) error {
	return pgx.BeginTxFunc(ctx, p.Pool, txOptions, func(tx pgx.Tx) error {
		return f(NewPgxTx(tx))
	})
}

// Get fetches a single row from the database and stores the result in the given struct.
func (p *PgxPool) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p, dst, query, args...)
}

// Select fetches multiple rows from the database and stores the results in the given slice of structs.
func (p *PgxPool) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.Pool, dest, query, args...)
}

// PgxTx is a PostgreSQL transaction wrapper that implements the Executor interface.
type PgxTx struct {
	pgx.Tx
}

func NewPgxTx(tx pgx.Tx) *PgxTx {
	return &PgxTx{tx}
}

func (p *PgxTx) Begin(ctx context.Context) (*PgxTx, error) {
	tx, err := p.Tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &PgxTx{tx}, nil
}

// BeginFunc starts a nested transaction and executes the given function within that transaction.
func (p *PgxTx) BeginFunc(ctx context.Context, f func(conn Tx) error) error {
	return pgx.BeginFunc(ctx, p.Tx, func(tx pgx.Tx) error {
		return f(NewPgxTx(tx))
	})
}

// Get fetches a single row from the database within the transaction and stores the result in the given struct.
func (p *PgxTx) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.Tx, dst, query, args...)
}

// Select fetches multiple rows from the database within the transaction and stores the results in the given slice of structs.
func (p *PgxTx) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.Tx, dest, query, args...)
}

// Config represents the configuration for establishing a connection pool.
type Config struct {
	DataSourceName  string        // PostgreSQL connection string
	MaxIdleConns    int           // Maximum number of idle connections in the pool
	MaxOpenConns    int           // Maximum number of open connections in the pool
	ConnMaxLifetime time.Duration // Maximum amount of time a connection can be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time a connection can be idle
	QueryLimit      uint          // Maximum number of rows to fetch in a single query
}

// NewConnection creates a new connection pool with the given configurationand returns a pointer to the pool.
func NewConnection(config Config) (*pgxpool.Pool, error) {
	parseConfig, err := pgxpool.ParseConfig(config.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to configure database, aborting: %w", err)
	}

	parseConfig.MaxConnLifetime = config.ConnMaxLifetime
	parseConfig.MaxConnIdleTime = config.ConnMaxIdleTime
	parseConfig.MaxConns = int32(config.MaxOpenConns)
	parseConfig.MinConns = int32(config.MaxIdleConns)

	pool, err := pgxpool.NewWithConfig(context.Background(), parseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database, aborting: %w", err)
	}

	// Ping the database to ensure a successful connection
	// Retry if unsuccessful for a limited number of attempts
	for i := 0; i <= pingMaxAttempts; i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutSecs*time.Second)
			defer cancel()

			return pool.Ping(ctx)
		}()

		if err == nil {
			break
		}

		if i < pingMaxAttempts {
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
	default:
		return false
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
		default:
			return false
		}

		return false
	}
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --name Querier
//go:generate go run github.com/vektra/mockery/v2@v2.42.0 --srcpkg github.com/jackc/pgx/v5 --inpackage=false --output . --outpkg=sql --filename=rows_mock.go --name Rows --structname MockRows
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

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
	pingMaxAttempts             = 6
	pingTimeoutSecs             = 10
	postgresUniqueViolationCode = "23505"
)

// compile time validator for the interfaces
var (
	_ Executor = PgxWrapper{}
	_ Executor = PgxTxWrapper{}

	ErrNoRows = pgx.ErrNoRows
)

// PgxWrapper is a PostgreSQL wrapper that implements the Executor interface.
type PgxWrapper struct {
	pool *pgxpool.Pool
}

// BeginFunc starts a transaction and executes the given function within that transaction.
func (p PgxWrapper) BeginFunc(ctx context.Context, f func(conn Executor) error) error {
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		return f(newPgxTxWrapper(tx))
	})
}

// Begin starts a new transaction.
func (p PgxWrapper) Begin(ctx context.Context) (*PgxTxWrapper, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPgxTxWrapper(tx), nil
}

// Exec executes a SQL statement and returns the result.
func (p PgxWrapper) Exec(ctx context.Context, query string, arguments ...any) (Result, error) {
	tag, err := p.pool.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

// Get fetches a single row from the database and stores the result in the given struct.
func (p PgxWrapper) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.pool, dst, query, args...)
}

// Select fetches multiple rows from the database and stores the results in the given slice of structs.
func (p PgxWrapper) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.pool, dest, query, args...)
}

// Close closes the connection pool.
func (p PgxWrapper) Close() {
	p.pool.Close()
}

// Ping checks the connection to the database.
func (p PgxWrapper) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

// PgxTxWrapper is a PostgreSQL transaction wrapper that implements the Executor interface.
type PgxTxWrapper struct {
	tx pgx.Tx
}

// BeginFunc starts a nested transaction and executes the given function within that transaction.
func (p PgxTxWrapper) BeginFunc(ctx context.Context, f func(conn Executor) error) error {
	return pgx.BeginFunc(ctx, p.tx, func(tx pgx.Tx) error {
		return f(PgxTxWrapper{tx})
	})
}

// Begin starts a new nested transaction.
func (p PgxTxWrapper) Begin(ctx context.Context) (*PgxTxWrapper, error) {
	tx, err := p.tx.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return newPgxTxWrapper(tx), nil
}

// Commit commits the transaction.
func (p PgxTxWrapper) Commit(ctx context.Context) error {
	return p.tx.Commit(ctx)
}

// Rollback rolls back the transaction.
func (p PgxTxWrapper) Rollback(ctx context.Context) error {
	return p.tx.Rollback(ctx)
}

// Exec executes a SQL statement within the transaction and returns the result.
func (p PgxTxWrapper) Exec(ctx context.Context, query string, arguments ...any) (Result, error) {
	tag, err := p.tx.Exec(ctx, query, arguments...)
	if err != nil {
		return nil, err
	}
	return pgxWrapperResult{tag}, nil
}

// Get fetches a single row from the database within the transaction and stores the result in the given struct.
func (p PgxTxWrapper) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, p.tx, dst, query, args...)
}

// Select fetches multiple rows from the database within the transaction and stores the results in the given slice of structs.
func (p PgxTxWrapper) Select(ctx context.Context, dest any, query string, args ...any) error {
	return pgxscan.Select(ctx, p.tx, dest, query, args...)
}

// pgxWrapperResult is a wrapper around pgconn.CommandTag that implements the Result interface.
type pgxWrapperResult struct {
	pgconn.CommandTag
}

// LastInsertId returns the ID of the last inserted row.
// Note: Not implemented in this wrapper.
func (r pgxWrapperResult) LastInsertId() (int64, error) {
	return 0, errors.New("not implemented")
}

// RowsAffected returns the number of rows affected by the operation.
func (r pgxWrapperResult) RowsAffected() (int64, error) {
	return r.CommandTag.RowsAffected(), nil
}

// NewPgxWrapper creates a new PgxWrapper with the given connection pool.
func NewPgxWrapper(pool *pgxpool.Pool) *PgxWrapper {
	return &PgxWrapper{pool}
}

// newPgxTxWrapper creates a new PgxTxWrapper with the given transaction.
func newPgxTxWrapper(tx pgx.Tx) *PgxTxWrapper {
	return &PgxTxWrapper{tx}
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

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name Querier
//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --srcpkg github.com/jackc/pgx/v5 --inpackage=false --output . --outpkg=sql --filename=rows_mock.go --name Rows --structname MockRows
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

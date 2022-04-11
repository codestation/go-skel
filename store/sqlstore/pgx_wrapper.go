package sqlstore

import (
	"context"
	"database/sql"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// compile time validator for the interfaces
var (
	_ sqlExecutor = pgxWrapper{}
	_ sqlExecutor = pgxTxWrapper{}

	ErrNoRows = pgx.ErrNoRows
)

type sqlExecutor interface {
	Begin(ctx context.Context, f func(db sqlExecutor) error) error
	Exec(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error)
	Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type pgxWrapper struct {
	*pgxpool.Pool
}

func (p pgxWrapper) Begin(ctx context.Context, f func(db sqlExecutor) error) error {
	return p.Pool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(newPgxTxWrapper(tx))
	})
}

func (p pgxWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error) {
	tag, err := p.Pool.Exec(ctx, sql, arguments...)
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

func (p pgxTxWrapper) Begin(ctx context.Context, f func(db sqlExecutor) error) error {
	return p.Tx.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(pgxTxWrapper{tx})
	})
}

func (p pgxTxWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) (sql.Result, error) {
	tag, err := p.Tx.Exec(ctx, sql, arguments...)
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

type sqlFuncExecutor interface {
	Close()
	Ping(ctx context.Context) error
}

type sqlDb interface {
	sqlExecutor
	sqlFuncExecutor
}

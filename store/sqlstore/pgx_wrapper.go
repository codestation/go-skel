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

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sql

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Executor interface {
	Begin(ctx context.Context) (*PgxTx, error)
	BeginFunc(ctx context.Context, f func(conn Tx) error) error
	Exec(ctx context.Context, query string, arguments ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, dst any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Transactor interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type Tx interface {
	Executor
	Transactor
}

type Connector interface {
	Executor
	Pinger
}

type Database interface {
	Connector
	Close()
}

type Listener interface {
	WaitForNotification(ctx context.Context) (*pgconn.Notification, error)
	Release()
}

type Acquirer interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
}

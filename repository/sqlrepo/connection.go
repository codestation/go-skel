// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlrepo

import (
	"context"
)

type SqlExecutor interface {
	Begin(ctx context.Context) (*PgxTxWrapper, error)
	BeginFunc(ctx context.Context, f func(conn SqlExecutor) error) error
	Exec(ctx context.Context, query string, arguments ...any) (Result, error)
	Get(ctx context.Context, dst any, query string, args ...any) error
	Select(ctx context.Context, dest any, query string, args ...any) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type SqlTransactor interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type SqlPinger interface {
	Ping(ctx context.Context) error
}

type SqlTx interface {
	SqlExecutor
	SqlTransactor
}

type SqlConnector interface {
	SqlExecutor
	SqlPinger
}

type SqlDb interface {
	SqlConnector
	Close()
}

// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"database/sql"
)

type SqlExecutor interface {
	BeginFunc(ctx context.Context, f func(db SqlExecutor) error) error
	Begin(ctx context.Context) (*PgxTxWrapper, error)
	Exec(ctx context.Context, query string, arguments ...interface{}) (sql.Result, error)
	Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type SqlFuncExecutor interface {
	Close()
	Ping(ctx context.Context) error
}

type SqlDb interface {
	SqlExecutor
	SqlFuncExecutor
}

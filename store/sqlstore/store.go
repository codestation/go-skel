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
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store"
)

type Stores struct {
	healthCheck store.HealthCheckStore
	// define more stores here
}

type SqlStore struct {
	db       SqlExecutor
	dbx      SqlFuncExecutor
	stores   Stores
	settings *model.SqlSettings
	builder  goqu.DialectWrapper
}

func (ss *SqlStore) initialize() {
	ss.builder = NewQueryBuilder()
	// Create all the stores here
	ss.stores.healthCheck = newSqlHealthCheckStore(ss)
}

func New(conn SqlDb, settings model.SqlSettings) *SqlStore {
	sqlStore := &SqlStore{
		settings: &settings,
	}

	sqlStore.db = conn
	sqlStore.dbx = conn

	// Create all the stores here
	sqlStore.initialize()

	return sqlStore
}

func (ss *SqlStore) HealthCheck() store.HealthCheckStore {
	return ss.stores.healthCheck
}

func (ss *SqlStore) Close() {}

func (ss *SqlStore) WithTransaction(ctx context.Context, f func(s store.Store) error) error {
	err := ss.db.BeginFunc(ctx, func(db SqlExecutor) error {
		s := &SqlStore{
			db:       db,
			settings: ss.settings,
		}
		s.initialize()
		return f(s)
	})

	if err != nil {
		return fmt.Errorf("sqlstore: transaction failed: %w", err)
	}
	return nil
}

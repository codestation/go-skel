// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"megpoid.xyz/go/go-skel/config"
	"megpoid.xyz/go/go-skel/store"
)

type Stores struct {
	healthCheck store.HealthCheckStore
	// define more stores here
	profile store.ProfileStore
}

type SqlStore struct {
	db       SqlExecutor
	dbx      SqlFuncExecutor
	stores   Stores
	settings *config.SqlSettings
	builder  goqu.DialectWrapper
}

func (ss *SqlStore) initialize() {
	ss.builder = NewQueryBuilder()
	// Create all the stores here
	ss.stores.healthCheck = newSqlHealthCheckStore(ss)
	ss.stores.profile = newSqlProfileStore(ss)
}

func New(conn SqlDb, settings config.SqlSettings) *SqlStore {
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

func (ss *SqlStore) Profile() store.ProfileStore {
	return ss.stores.profile
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

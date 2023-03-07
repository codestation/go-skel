// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/repository/sqlrepo"
	"megpoid.dev/go/go-skel/test"
)

type Connection struct {
	db              sqlrepo.SqlDb
	store           sqlrepo.SqlExecutor
	tx              sqlrepo.SqlTransactor
	withTransaction bool
}

func NewTestConnection(t *testing.T, withTransaction bool) *Connection {
	c := &Connection{withTransaction: withTransaction}
	c.setupDatabase(t)
	return c
}

func (c *Connection) Close(t *testing.T) {
	t.Helper()
	if c.withTransaction {
		// do not save the changes to the database
		err := c.tx.Rollback(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to rollback transaction", err.Error())
		}
	}
	c.db.Close()
}

func (c *Connection) setupDatabase(t *testing.T) {
	t.Helper()

	cfg, err := config.NewConfig()
	if err != nil {
		assert.FailNowf(t, "Cannot load settings", err.Error())
	}

	if dsn, ok := os.LookupEnv("APP_DSN"); ok {
		cfg.SqlSettings.DataSourceName = dsn
	}

	conn, err := sqlrepo.NewConnection(cfg.SqlSettings)
	if err != nil {
		assert.FailNowf(t, "Failed to create database Connection", err.Error())
	}

	c.db = sqlrepo.NewPgxWrapper(conn)

	if c.withTransaction {
		// create a new transaction so a test doesn't interfere with another
		tx, err := c.db.Begin(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to create transaction", err.Error())
		}

		c.seedDatabase(t, tx)
		c.tx = tx
		c.store = tx
	} else {
		c.store = c.db
	}
}

func (c *Connection) seedDatabase(t *testing.T, conn sqlrepo.SqlExecutor) {
	assets := test.Assets()
	data, err := assets.ReadFile("seed/base.sql")
	if err != nil {
		assert.FailNowf(t, "Failed to read seed file", err.Error())
	}

	_, err = conn.Exec(context.Background(), string(data))
	if err != nil {
		assert.FailNowf(t, "Failed to run seed file", err.Error())
	}
}

type fakeDatabase struct {
	Error  error
	Result *fakeSqlResult
}

func (d fakeDatabase) BeginFunc(ctx context.Context, f func(conn sqlrepo.SqlExecutor) error) error {
	return d.Error
}

func (d fakeDatabase) Begin(ctx context.Context) (*sqlrepo.PgxTxWrapper, error) {
	return nil, d.Error
}

func (d fakeDatabase) Exec(ctx context.Context, query string, arguments ...any) (sql.Result, error) {
	if d.Result != nil {
		return d.Result, nil
	}
	return nil, d.Error
}

func (d fakeDatabase) Get(ctx context.Context, dst any, query string, args ...any) error {
	return d.Error
}

func (d fakeDatabase) Select(ctx context.Context, dest any, query string, args ...any) error {
	return d.Error
}

var _ sqlrepo.SqlExecutor = &fakeDatabase{}

type fakeSqlResult struct {
	Error error
}

func (f fakeSqlResult) LastInsertId() (int64, error) {
	return 0, f.Error
}

func (f fakeSqlResult) RowsAffected() (int64, error) {
	return 0, f.Error
}

var _ sql.Result = &fakeSqlResult{}

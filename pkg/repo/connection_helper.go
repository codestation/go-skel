// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/pkg/sql"
	"megpoid.dev/go/go-skel/testdata"
)

type Connection struct {
	Db              sql.Database
	Store           sql.Executor
	tx              sql.Transactor
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
	c.Db.Close()
}

func (c *Connection) setupDatabase(t *testing.T) {
	t.Helper()

	cfg := sql.Config{
		DataSourceName:  "postgres://goapp:secret@localhost/goapp?sslmode=disable",
		ConnMaxLifetime: 1 * time.Hour,
		ConnMaxIdleTime: 5 * time.Minute,
		MaxOpenConns:    10,
		MaxIdleConns:    1,
	}

	if dsn, ok := os.LookupEnv("APP_DSN"); ok {
		cfg.DataSourceName = dsn
	}

	conn, err := sql.NewConnection(cfg)
	if err != nil {
		assert.FailNowf(t, "Failed to create database Connection", err.Error())
	}

	c.Db = sql.NewPgxWrapper(conn)

	if c.withTransaction {
		// create a new transaction so a test doesn't interfere with another
		tx, err := c.Db.Begin(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to create transaction", err.Error())
		}

		c.seedDatabase(t, tx)
		c.tx = tx
		c.Store = tx
	} else {
		c.Store = c.Db
	}
}

func (c *Connection) seedDatabase(t *testing.T, conn sql.Executor) {
	assets := testdata.SqlAssets()
	data, err := assets.ReadFile("sql/testdata.sql")
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
	Result *fakeSQLResult
}

func (d fakeDatabase) BeginFunc(ctx context.Context, f func(conn sql.Executor) error) error {
	return d.Error
}

func (d fakeDatabase) Begin(ctx context.Context) (*sql.PgxTxWrapper, error) {
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

var _ sql.Executor = &fakeDatabase{}

type fakeSQLResult struct {
	Error error
}

func (f fakeSQLResult) LastInsertId() (int64, error) {
	return 0, f.Error
}

func (f fakeSQLResult) RowsAffected() (int64, error) {
	return 0, f.Error
}

var _ sql.Result = &fakeSQLResult{}

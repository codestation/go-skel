// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/config"
	"megpoid.xyz/go/go-skel/testdata"
	"os"
	"testing"
)

type connection struct {
	store           *SqlStore
	db              *PgxWrapper
	withTransaction bool
}

func NewTestConnection(t *testing.T, withTransaction bool) *connection {
	c := &connection{withTransaction: withTransaction}
	c.setupDatabase(t)
	return c
}

func (c *connection) Close(t *testing.T) {
	t.Helper()
	if c.withTransaction {
		// do not save the changes to the database
		err := c.store.Rollback(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to rollback transaction", err.Error())
		}
	}
	c.db.Close()
}

func (c *connection) setupDatabase(t *testing.T) {
	t.Helper()

	dsn := os.Getenv("APP_DSN")
	if dsn == "" {
		t.Skip("Test skipped because APP_DSN environment variable wasn't found")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		assert.FailNowf(t, "Cannot load settings", err.Error())
	}

	cfg.SqlSettings.DataSourceName = dsn

	conn, err := NewConnection(cfg.SqlSettings)
	if err != nil {
		assert.FailNowf(t, "Failed to create database connection", err.Error())
	}

	c.db = conn

	store := New(conn, cfg.SqlSettings)

	if c.withTransaction {
		// create a new transaction so a test doesn't interfere with another
		tx, err := store.Begin(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to create transaction", err.Error())
		}

		c.seedDatabase(t, tx.db)
		c.store = tx
	} else {
		c.store = store
	}
}

func (c *connection) seedDatabase(t *testing.T, conn SqlExecutor) {
	assets := testdata.Assets()
	data, err := assets.ReadFile("seed/base.sql")
	if err != nil {
		assert.FailNowf(t, "Failed to read seed file", err.Error())
	}

	_, err = conn.Exec(context.Background(), string(data))
	if err != nil {
		assert.FailNowf(t, "Failed to run seed file", err.Error())
	}
}

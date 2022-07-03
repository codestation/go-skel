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

func setupDatabase(t *testing.T) *SqlStore {
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
	store := New(conn, cfg.SqlSettings)
	// create a new transaction so a test doesn't interfere with another
	tx, err := store.Begin(context.Background())
	if err != nil {
		assert.FailNowf(t, "Failed to create transaction", err.Error())
	}

	seedDatabase(t, tx.db)

	return tx
}

func seedDatabase(t *testing.T, conn SqlExecutor) {
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

func teardownDatabase(t *testing.T, store *SqlStore) {
	t.Helper()
	if store != nil {
		// do not save the changes to the database
		err := store.Rollback(context.Background())
		if err != nil {
			assert.FailNowf(t, "Failed to rollback transaction", err.Error())
		}
		store.Close()
	}
}

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
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store"
)

type FakeDbConn struct{}

func (d *FakeDbConn) Begin(_ context.Context, f func(db SqlExecutor) error) error {
	db := &FakeDbConn{}
	return f(db)
}

func (d *FakeDbConn) Exec(_ context.Context, _ string, x ...interface{}) (sql.Result, error) {
	panic("implement me")
}

func (d *FakeDbConn) Get(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	panic("implement me")
}

func (d *FakeDbConn) Select(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	panic("implement me")
}

func (d *FakeDbConn) Close() {}

func (d *FakeDbConn) Ping(_ context.Context) error {
	return nil
}

func TestNew(t *testing.T) {
	db := &FakeDbConn{}
	ss := New(db, model.SqlSettings{})
	assert.NotNil(t, ss.stores.healthCheck)
	assert.NotNil(t, ss.db)
	assert.NotNil(t, ss.dbx)
	assert.NotNil(t, ss.settings)
	assert.NotNil(t, ss.builder)
}

func TestSqlStore_WithTransaction(t *testing.T) {
	db := &FakeDbConn{}
	ss := &SqlStore{db: db}
	err := ss.WithTransaction(context.Background(), func(s store.Store) error {
		tx, ok := s.(*SqlStore)
		assert.True(t, ok)
		assert.NotEqual(t, ss, tx)
		return nil
	})
	assert.NoError(t, err)
}

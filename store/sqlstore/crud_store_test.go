// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"github.com/stretchr/testify/suite"
	"megpoid.xyz/go/go-skel/model"
	"testing"
	"time"
)

type testUser struct {
	model.Model
	Name string
}

func TestStore(t *testing.T) {
	suite.Run(t, &crudSuite{})
}

type crudSuite struct {
	suite.Suite
	conn *connection
}

func (s *crudSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), true)
}

func (s *crudSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *crudSuite) TestNewStore() {
	store := NewCrudStore[model.Profile](s.conn.store)
	s.Equal("profiles", store.table)
	s.Equal([]any{"*"}, store.selectFields)
}

func (s *crudSuite) TestStoreGet() {
	store := NewCrudStore[testUser](s.conn.store)
	user, err := store.Get(context.Background(), 1)
	if s.NoError(err) {
		s.Equal(model.ID(1), user.ID)
	}
}

func (s *crudSuite) TestStoreList() {
	store := NewCrudStore[testUser](s.conn.store)
	users, err := store.List(context.Background())
	if s.NoError(err) {
		s.Len(users.Data, 1)
	}
}

func (s *crudSuite) TestStoreUpdate() {
	store := NewCrudStore[testUser](s.conn.store)
	user, err := store.Get(context.Background(), 1)
	if s.NoError(err) {
		user.Name = "Jane Doe"
		user.UpdatedAt = time.Now()
		err := store.Update(context.Background(), user)
		s.NoError(err)
	}
}

func (s *crudSuite) TestStoreRemove() {
	store := NewCrudStore[testUser](s.conn.store)
	err := store.Delete(context.Background(), 1)
	s.NoError(err)
}

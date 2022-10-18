// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"github.com/stretchr/testify/suite"
	"megpoid.dev/go/go-skel/store/clause"
	"megpoid.dev/go/go-skel/store/filter"
	"testing"
)

func TestProfileStore(t *testing.T) {
	suite.Run(t, &profileSuite{})
}

type profileSuite struct {
	suite.Suite
	conn *Connection
}

func (s *profileSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), true)
}

func (s *profileSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *profileSuite) TestFilterSingleMatch() {
	store := newSqlProfileStore(s.conn.store)
	result, err := store.List(context.Background(), clause.WithConditions(filter.Condition{
		Field:     "first_name",
		Operation: filter.OperationEqual,
		Value:     "John",
	}))
	if s.NoError(err) {
		s.Equal(1, len(result.Data))
	}
}

func (s *profileSuite) TestFilterMultipleMatch() {
	store := newSqlProfileStore(s.conn.store)
	result, err := store.List(context.Background(), clause.WithConditions(filter.Condition{
		Field:     "last_name",
		Operation: filter.OperationEqual,
		Value:     "Doe",
	}))
	if s.NoError(err) {
		s.Equal(2, len(result.Data))
	}
}

func (s *profileSuite) TestFilterNoMatch() {
	store := newSqlProfileStore(s.conn.store)
	result, err := store.List(context.Background(), clause.WithConditions(filter.Condition{
		Field:     "last_name",
		Operation: filter.OperationEqual,
		Value:     "Unknown",
	}))
	if s.NoError(err) {
		s.Equal(0, len(result.Data))
	}
}

func (s *profileSuite) TestProfileByUserToken() {
	store := newSqlProfileStore(s.conn.store)
	profile, err := store.GetByUserToken(context.Background(), "1234")
	if s.NoError(err) {
		s.Equal("1234", profile.UserToken)
	}
}

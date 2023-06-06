// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"megpoid.dev/go/go-skel/pkg/clause"
	"megpoid.dev/go/go-skel/pkg/repo"
	filter2 "megpoid.dev/go/go-skel/pkg/repo/filter"
)

func TestProfileStore(t *testing.T) {
	suite.Run(t, &profileSuite{})
}

type profileSuite struct {
	suite.Suite
	conn *repo.Connection
}

func (s *profileSuite) SetupTest() {
	s.conn = repo.NewTestConnection(s.T(), true)
}

func (s *profileSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *profileSuite) TestFilterSingleMatch() {
	store := NewProfileRepo(s.conn.Store)
	result, err := store.List(context.Background(), clause.WithConditions(filter2.Condition{
		Field:     "first_name",
		Operation: filter2.OperationEqual,
		Value:     "John",
	}))
	if s.NoError(err) {
		s.Equal(1, len(result.Items))
	}
}

func (s *profileSuite) TestFilterMultipleMatch() {
	store := NewProfileRepo(s.conn.Store)
	result, err := store.List(context.Background(), clause.WithConditions(filter2.Condition{
		Field:     "last_name",
		Operation: filter2.OperationEqual,
		Value:     "Doe",
	}))
	if s.NoError(err) {
		s.Equal(2, len(result.Items))
	}
}

func (s *profileSuite) TestFilterNoMatch() {
	store := NewProfileRepo(s.conn.Store)
	result, err := store.List(context.Background(), clause.WithConditions(filter2.Condition{
		Field:     "last_name",
		Operation: filter2.OperationEqual,
		Value:     "Unknown",
	}))
	if s.NoError(err) {
		s.Equal(0, len(result.Items))
	}
}

func (s *profileSuite) TestProfileByEmail() {
	store := NewProfileRepo(s.conn.Store)
	profile, err := store.GetByEmail(context.Background(), "john.doe@example.com")
	if s.NoError(err) {
		s.Equal("john.doe@example.com", profile.Email)
	}
}

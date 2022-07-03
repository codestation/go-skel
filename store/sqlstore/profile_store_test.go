// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"github.com/stretchr/testify/suite"
	"megpoid.xyz/go/go-skel/model"
	"testing"
)

func TestProfileStore(t *testing.T) {
	suite.Run(t, &profileSuite{})
}

type profileSuite struct {
	suite.Suite
	sqlStore *SqlStore
}

func (s *profileSuite) SetupTest() {
	s.sqlStore = setupDatabase(s.T())
}

func (s *profileSuite) TearDownTest() {
	teardownDatabase(s.T(), s.sqlStore)
}

func (s *profileSuite) TestNewStore() {
	store := NewCrudStore[model.Profile](s.sqlStore)
	s.Equal("profiles", store.table)
	s.Equal([]any{"*"}, store.selectFields)
}

func (s *profileSuite) TestProfileByUserToken() {
	store := newSqlProfileStore(s.sqlStore)
	profile, err := store.GetByUserToken(context.Background(), "1234")
	if s.NoError(err) {
		s.Equal("1234", profile.UserToken)
	}
}

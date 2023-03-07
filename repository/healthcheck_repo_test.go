// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHealthcheckStore(t *testing.T) {
	suite.Run(t, &healthcheckSuite{})
}

type healthcheckSuite struct {
	suite.Suite
	conn *Connection
}

func (s *healthcheckSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), false)
}

func (s *healthcheckSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *healthcheckSuite) TestPing() {
	store := NewHealthcheckRepo(s.conn.db)
	err := store.Execute(context.Background())
	s.NoError(err)
}

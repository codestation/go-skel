// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package uow

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"megpoid.dev/go/go-skel/pkg/repo"
)

func TestUnitOfWork(t *testing.T) {
	suite.Run(t, &unitOfWorkSuite{})
}

type unitOfWorkSuite struct {
	suite.Suite
	conn *repo.Connection
}

func (s *unitOfWorkSuite) SetupTest() {
	s.conn = repo.NewTestConnection(s.T(), false)
}

func (s *unitOfWorkSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *unitOfWorkSuite) TestTransaction() {
	uow := New(s.conn.Db)
	err := uow.Do(context.Background(), func(work UnitOfWork) error {
		return nil
	})

	assert.NoError(s.T(), err)
}

func (s *unitOfWorkSuite) TestTransactionNested() {
	uow := New(s.conn.Db)
	err1 := uow.Do(context.Background(), func(tx1 UnitOfWork) error {
		err2 := tx1.Do(context.Background(), func(tx2 UnitOfWork) error {
			return nil
		})

		return err2
	})

	assert.NoError(s.T(), err1)
}

func (s *unitOfWorkSuite) TestTransactionRollback() {
	uow := New(s.conn.Db)
	err := uow.Do(context.Background(), func(work UnitOfWork) error {
		return errors.New("an error")
	})

	assert.Error(s.T(), err)
}

func (s *unitOfWorkSuite) TestTransactionNestedRollback() {
	uow := New(s.conn.Db)
	myErr := errors.New("an error")
	err1 := uow.Do(context.Background(), func(tx1 UnitOfWork) error {
		err2 := tx1.Do(context.Background(), func(tx2 UnitOfWork) error {
			return myErr
		})

		// ignore if custom error
		if !errors.Is(err2, myErr) {
			return err2
		}

		return nil
	})

	assert.NoError(s.T(), err1)
}

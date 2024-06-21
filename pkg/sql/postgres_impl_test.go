// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sql

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsUniqueError(t *testing.T) {
	err := &pgconn.PgError{Code: postgresUniqueViolationCode}
	uniq := IsUniqueError(err)
	assert.True(t, uniq)
}

func TestIsUniqueErrorFail(t *testing.T) {
	err := &pgconn.PgError{Code: "0"}
	uniq := IsUniqueError(err)
	assert.False(t, uniq)
}

func TestIsUniqueErrorOther(t *testing.T) {
	err := errors.New("an error")
	uniq := IsUniqueError(err)
	assert.False(t, uniq)
}

func TestIsUniqueErrorConstraint(t *testing.T) {
	err := &pgconn.PgError{Code: postgresUniqueViolationCode, ConstraintName: "foo"}
	uniq := IsUniqueError(err, WithConstraintName("foo"))
	assert.True(t, uniq)
}

func TestIsUniqueErrorConstraintFalse(t *testing.T) {
	err := &pgconn.PgError{Code: postgresUniqueViolationCode, ConstraintName: "foo"}
	uniq := IsUniqueError(err, WithConstraintName("bar"))
	assert.False(t, uniq)
}

func TestIsUniqueErrorConstraintOther(t *testing.T) {
	err := errors.New("an error")
	uniq := IsUniqueError(err, WithConstraintName("bar"))
	assert.False(t, uniq)
}

type Foo struct {
	ID   int
	Name string
}

func TestGetStruct(t *testing.T) {
	db := NewMockQuerier(t)
	rows := NewMockRows(t)
	fields := []pgconn.FieldDescription{
		{Name: "id"},
		{Name: "name"},
	}
	db.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything).Return(rows, nil)
	rows.EXPECT().Next().Return(true)
	rows.EXPECT().FieldDescriptions().Return(fields)
	rows.EXPECT().Scan(mock.Anything, mock.Anything).RunAndReturn(func(i ...any) error {
		if len(i) != 2 {
			return errors.New("scan: invalid arg count")
		}
		idVal, ok := i[0].(*int)
		if !ok {
			return errors.New("scan: invalid scan type")
		}

		nameVal, ok := i[1].(*string)
		if !ok {
			return errors.New("scan: invalid scan type")
		}

		*idVal = 1
		*nameVal = "John"

		return nil
	})

	rows.EXPECT().Close()
	rows.EXPECT().Err().Return(nil)

	result, err := GetStruct[Foo](context.Background(), db, "")
	assert.NoError(t, err)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "John", result.Name)
}

func TestSelectStruct(t *testing.T) {
	db := NewMockQuerier(t)
	rows := NewMockRows(t)
	fields := []pgconn.FieldDescription{
		{Name: "id"},
		{Name: "name"},
	}
	db.EXPECT().Query(mock.Anything, mock.Anything, mock.Anything).Return(rows, nil)
	rows.EXPECT().Next().Return(true).Times(3)
	rows.EXPECT().Next().Return(false)
	rows.EXPECT().FieldDescriptions().Return(fields)
	rows.EXPECT().Scan(mock.Anything, mock.Anything).RunAndReturn(func(i ...any) error {
		if len(i) != 2 {
			return errors.New("scan: invalid arg count")
		}
		idVal, ok := i[0].(*int)
		if !ok {
			return errors.New("scan: invalid scan type")
		}

		nameVal, ok := i[1].(*string)
		if !ok {
			return errors.New("scan: invalid scan type")
		}

		*idVal = 1
		*nameVal = "John"

		return nil
	})
	rows.EXPECT().Close()
	rows.EXPECT().Err().Return(nil)

	results, err := SelectStruct[Foo](context.Background(), db, "")
	assert.NoError(t, err)
	assert.Len(t, results, 3)
	assert.Equal(t, 1, results[0].ID)
	assert.Equal(t, "John", results[0].Name)
}

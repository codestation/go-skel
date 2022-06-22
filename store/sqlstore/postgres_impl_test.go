// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"errors"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
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

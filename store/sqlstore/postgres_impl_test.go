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

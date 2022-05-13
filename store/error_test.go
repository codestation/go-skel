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

package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var err = errors.New("test error")
var internalErr = errors.New("internal error")

func TestNewRepoError(t *testing.T) {
	t.Run("full message", func(t *testing.T) {
		repoErr := NewRepoError(err, internalErr)
		assert.Equal(t, repoErr.Error(), err.Error()+": "+internalErr.Error())
	})
	t.Run("no internal", func(t *testing.T) {
		repoErr := NewRepoError(err, nil)
		assert.Equal(t, repoErr.Error(), err.Error())
	})
}

func TestRepoError_Unwrap(t *testing.T) {
	repoErr := NewRepoError(err, internalErr)
	assert.ErrorIs(t, repoErr, err)
}

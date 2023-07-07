// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsPointerInt(t *testing.T) {
	x := 1
	y := AsPointer(x)
	if assert.NotNil(t, y) {
		assert.Equal(t, 1, *y)
	}
}

func TestAsPointerString(t *testing.T) {
	x := "test"
	y := AsPointer(x)
	if assert.NotNil(t, y) {
		assert.Equal(t, "test", *y)
	}
}

func TestAsPointerMap(t *testing.T) {
	x := map[int]string{1: "test"}
	y := AsPointer(x)
	if assert.NotNil(t, y) {
		assert.Equal(t, map[int]string{1: "test"}, *y)
	}
}

func TestAsValue(t *testing.T) {
	x := new(int)
	*x = 1
	y := AsValue(x)
	assert.Equal(t, 1, y)
}

func TestAsValueNil(t *testing.T) {
	var x *string
	y := AsValue(x)
	assert.Equal(t, "", y)
}

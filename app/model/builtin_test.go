// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewType(t *testing.T) {
	intVal := 1
	value := NewType(intVal)
	var intPointer *int
	assert.IsType(t, intPointer, value)
	assert.Equal(t, 1, *value)
	assert.NotEqual(t, &intVal, intPointer)
}

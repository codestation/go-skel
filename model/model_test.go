// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	now := time.Now()
	model := NewModel(now)
	assert.Equal(t, model.CreatedAt, now)
	assert.Equal(t, model.UpdatedAt, now)
}

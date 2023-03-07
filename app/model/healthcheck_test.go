// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckResult(t *testing.T) {
	t.Run("AllOk", func(t *testing.T) {
		healthCheck := HealthcheckResult{Ping: nil}
		assert.True(t, healthCheck.AllOk())
	})
	t.Run("Fail", func(t *testing.T) {
		healthCheck := HealthcheckResult{Ping: errors.New("an error")}
		assert.False(t, healthCheck.AllOk())
	})
}

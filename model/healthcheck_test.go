package model

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheckResult(t *testing.T) {
	t.Run("AllOk", func(t *testing.T) {
		healthCheck := HealthCheckResult{Ping: nil}
		assert.True(t, healthCheck.AllOk())
	})
	t.Run("Fail", func(t *testing.T) {
		healthCheck := HealthCheckResult{Ping: errors.New("an error")}
		assert.False(t, healthCheck.AllOk())
	})
}

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

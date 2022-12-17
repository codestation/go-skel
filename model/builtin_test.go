package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewType(t *testing.T) {
	value := NewType(1)
	var intPointer *int
	assert.IsType(t, intPointer, value)
	assert.Equal(t, 1, *value)
}

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsPointerInt(t *testing.T) {
	var x = 1
	y := AsPointer(x)
	if assert.NotNil(t, y) {
		assert.Equal(t, 1, *y)
	}
}

func TestAsPointerString(t *testing.T) {
	var x = "test"
	y := AsPointer(x)
	if assert.NotNil(t, y) {
		assert.Equal(t, "test", *y)
	}
}

func TestAsPointerMap(t *testing.T) {
	var x = map[int]string{1: "test"}
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
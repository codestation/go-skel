package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssets(t *testing.T) {
	fs := Assets()
	entries, err := fs.ReadDir("static")
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)
}

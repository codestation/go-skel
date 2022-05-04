package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssets(t *testing.T) {
	fs := Assets()
	entries, err := fs.ReadDir("migrations")
	assert.NoError(t, err)
	assert.NotEmpty(t, entries)
}

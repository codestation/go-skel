package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var err = errors.New("test error")
var internalErr = errors.New("internal error")

func TestNewRepoError(t *testing.T) {
	t.Run("full message", func(t *testing.T) {
		repoErr := NewRepoError(err, internalErr)
		assert.Equal(t, repoErr.Error(), err.Error()+": "+internalErr.Error())
	})
	t.Run("no internal", func(t *testing.T) {
		repoErr := NewRepoError(err, nil)
		assert.Equal(t, repoErr.Error(), err.Error())
	})
}

func TestRepoError_Unwrap(t *testing.T) {
	repoErr := NewRepoError(err, internalErr)
	assert.ErrorIs(t, repoErr, err)
}

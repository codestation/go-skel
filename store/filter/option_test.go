package filter

import (
	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/model/request"
	"testing"
)

func TestEmptyOptions(t *testing.T) {
	f := New()
	assert.Empty(t, f.rules)
	assert.Empty(t, f.filters)
}

func TestOptionsWith(t *testing.T) {
	opts := []Option{
		WithFilters([]request.Filter{
			{Field: "foo"},
		}...),
		WithRules([]Rule{
			{Key: "bar"},
		}...),
	}

	f := New(opts...)
	assert.Len(t, f.filters, 1)
	assert.Len(t, f.rules, 1)
	assert.Contains(t, f.rules, "bar")
}

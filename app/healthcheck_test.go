package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.xyz/go/go-skel/store"
)

type fakeStore struct{}

func (s fakeStore) HealthCheck() store.HealthCheckStore {
	return &fakeHealthCheckStore{}
}

func (s fakeStore) WithTransaction(_ context.Context, _ func(s store.Store) error) error {
	panic("implement me")
}

type fakeHealthCheckStore struct{}

func (f fakeHealthCheckStore) HealthCheck(_ context.Context) error {
	return nil
}

func TestApp_HealthCheck(t *testing.T) {
	ss := &fakeStore{}
	srv := &Server{Store: ss}
	app := New(srv)
	result := app.HealthCheck(context.Background())
	assert.NotNil(t, result)
	assert.True(t, result.AllOk())
}

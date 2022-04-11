package sqlstore

import (
	"context"

	"megpoid.xyz/go/go-skel/store"
)

type SqlHealthCheckStore struct {
	*SqlStore
}

func newSqlHealthCheckStore(sqlStore *SqlStore) store.HealthCheckStore {
	s := &SqlHealthCheckStore{
		SqlStore: sqlStore,
	}

	return s
}

// HealthCheck returns an error if the database doesn't respond
func (s SqlHealthCheckStore) HealthCheck(ctx context.Context) error {
	if err := s.dbx.Ping(ctx); err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}
	return nil
}

package sqlstore

import (
	"context"
	store2 "megpoid.xyz/go/go-skel/store"
)

type SqlHealthCheckStore struct {
	*SqlStore
}

func newSqlHealthCheckStore(sqlStore *SqlStore) store2.HealthCheckStore {
	s := &SqlHealthCheckStore{
		SqlStore: sqlStore,
	}

	return s
}

// HealthCheck returns an error if the database doesn't respond
func (s SqlHealthCheckStore) HealthCheck(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return store2.NewRepoError(store2.ErrBackend, err)
	}
	return nil
}

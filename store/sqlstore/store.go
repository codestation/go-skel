package sqlstore

import (
	"megpoid.xyz/go/go-skel/config"
	"megpoid.xyz/go/go-skel/store"
)

type Stores struct {
	healthCheck store.HealthCheckStore
	// define more stores here
}

type SqlStore struct {
	db     SQLConn
	stores Stores
}

func New(cfg *config.Config) *SqlStore {
	sqlStore := &SqlStore{}

	// Database initialization
	sqlStore.db = sqlStore.setupConnection(cfg)

	// Create all the stores here
	sqlStore.stores.healthCheck = newSqlHealthCheckStore(sqlStore)

	return sqlStore
}

func (ss *SqlStore) HealthCheck() store.HealthCheckStore {
	return ss.stores.healthCheck
}

func (ss *SqlStore) Close() error {
	return ss.db.Close()
}

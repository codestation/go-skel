package sqlstore

import (
	"context"
	"log"

	migrate "github.com/rubenv/sql-migrate"
	"megpoid.xyz/go/go-skel/db"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store"
)

type Stores struct {
	healthCheck store.HealthCheckStore
	// define more stores here
}

type SqlStore struct {
	db            SQLConn
	stores        Stores
	settings      *model.SqlSettings
	runMigrations bool
}

func New(settings model.SqlSettings) *SqlStore {
	sqlStore := &SqlStore{
		settings: &settings,
	}

	// Database initialization
	sqlStore.db = sqlStore.setupConnection()

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

func (ss *SqlStore) RunMigrations(settings model.MigrationSettings) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: db.Assets(),
		Root:       "migrations",
	}

	migrate.SetTable("app_migrations")
	ctx := context.Background()

	if settings.Reset {
		_, err := ss.db.ExecContext(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		if err != nil {
			return err
		}

		_, err = ss.db.ExecContext(ctx, "CREATE SCHEMA public")
		if err != nil {
			return nil
		}
		log.Printf("Recreated 'public' schema")
	}

	step := 0
	// SQLConn -> sqlx -> sql
	sqlDb := ss.db.DB().DB

	if !settings.Reset && (settings.Rollback || settings.Redo) {
		step = settings.Step
		n, err := migrate.ExecMax(sqlDb, ss.settings.DriverName, migrations, migrate.Down, step)
		if err != nil {
			return err
		}
		log.Printf("Reverted %d migrations", n)
	}

	if settings.Reset || !settings.Rollback || settings.Redo {
		if settings.Redo {
			step = settings.Step
		}

		n, err := migrate.ExecMax(sqlDb, ss.settings.DriverName, migrations, migrate.Up, step)
		if err != nil {
			return err
		}
		log.Printf("Applied %d migrations", n)
	}
	return nil
}

package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
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
	db       sqlExecutor
	dbx      sqlFuncExecutor
	stores   Stores
	settings *model.SqlSettings
}

func (ss *SqlStore) initialize() {
	// Create all the stores here
	ss.stores.healthCheck = newSqlHealthCheckStore(ss)
}

func New(settings model.SqlSettings) *SqlStore {
	sqlStore := &SqlStore{
		settings: &settings,
	}

	// Database initialization
	conn := sqlStore.setupConnection()
	sqlStore.db = conn
	sqlStore.dbx = conn

	// Create all the stores here
	sqlStore.initialize()

	return sqlStore
}

func (ss *SqlStore) HealthCheck() store.HealthCheckStore {
	return ss.stores.healthCheck
}

func (ss *SqlStore) Close() {
	ss.dbx.Close()
}

func (ss *SqlStore) RunMigrations(settings model.MigrationSettings) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: db.Assets(),
		Root:       "migrations",
	}

	migrate.SetTable("app_migrations")
	ctx := context.Background()

	if settings.Reset {
		_, err := ss.db.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		if err != nil {
			return err
		}

		_, err = ss.db.Exec(ctx, "CREATE SCHEMA public")
		if err != nil {
			return nil
		}
		log.Printf("Recreated 'public' schema")
	}

	step := 0
	//TODO: find a way to either convert pgx.Pool to sql.DD
	sqlDb, err := sql.Open("pgx", ss.settings.DataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open migration connection: %w", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to clone migration connection: %s", err.Error())
		}
	}(sqlDb)

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

func (ss *SqlStore) WithTransaction(ctx context.Context, f func(s store.Store) error) error {
	err := ss.db.Begin(ctx, func(db sqlExecutor) error {
		s := &SqlStore{
			db:       db,
			settings: ss.settings,
		}
		s.initialize()
		return f(s)
	})

	if err != nil {
		return fmt.Errorf("sqlstore: transaction failed: %w", err)
	}
	return nil
}

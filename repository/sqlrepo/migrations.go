// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlrepo

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"

	migrate "github.com/heroiclabs/sql-migrate"
	"github.com/jackc/pgx/v5/pgxpool"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/db"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, config *config.Config) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: db.Assets(),
		Root:       "migrations",
	}

	migration := migrate.MigrationSet{
		TableName: "app_migrations",
	}

	if config.MigrationSettings.Reset {
		_, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, "CREATE SCHEMA public")
		if err != nil {
			return nil
		}
		log.Printf("Recreated 'public' schema")
	}

	step := 0

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}

	defer conn.Release()

	if !config.MigrationSettings.Reset && (config.MigrationSettings.Rollback || config.MigrationSettings.Redo) {
		step = config.MigrationSettings.Step
		n, err := migration.ExecMax(ctx, conn.Conn(), migrations, migrate.Down, step)
		if err != nil {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		log.Printf("Reverted %d migrations", n)
	}

	if config.MigrationSettings.Reset || !config.MigrationSettings.Rollback || config.MigrationSettings.Redo {
		if config.MigrationSettings.Redo {
			step = config.MigrationSettings.Step
		}

		n, err := migration.ExecMax(ctx, conn.Conn(), migrations, migrate.Up, step)
		if err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		log.Printf("Applied %d migrations", n)
	}

	if config.MigrationSettings.Seed {
		if err := ApplySeed(ctx, pool); err != nil {
			return fmt.Errorf("failed to apply database seed: %w", err)
		}
	}

	return nil
}

// ApplySeed initializes the database with data
func ApplySeed(ctx context.Context, pool *pgxpool.Pool) error {
	assets := db.Seeds()
	entries, err := assets.ReadDir("seed")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			seed, err := assets.ReadFile(path.Join("seed", entry.Name()))
			if err != nil {
				return err
			}
			fmt.Printf("Applying seed to database: %s\n", entry.Name())
			if _, err := pool.Exec(ctx, string(seed)); err != nil {
				return err
			}
		}
	}

	return nil
}

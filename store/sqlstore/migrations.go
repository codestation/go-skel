// Copyright (c) 2022 codestation
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
	"log"
	"megpoid.xyz/go/go-skel/db"
	"megpoid.xyz/go/go-skel/model"
	"path"
	"strings"
)

func RunMigrations(ctx context.Context, pool SqlExecutor, config *model.Config) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: db.Assets(),
		Root:       "migrations",
	}

	migrate.SetTable("app_migrations")

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

	sqlPool, err := sql.Open("pgx", config.SqlSettings.DataSourceName)
	if err != nil {
		return err
	}

	defer sqlPool.Close()

	if !config.MigrationSettings.Reset && (config.MigrationSettings.Rollback || config.MigrationSettings.Redo) {
		step = config.MigrationSettings.Step
		n, err := migrate.ExecMax(sqlPool, config.SqlSettings.DriverName, migrations, migrate.Down, step)
		if err != nil {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		log.Printf("Reverted %d migrations", n)
	}

	if config.MigrationSettings.Reset || !config.MigrationSettings.Rollback || config.MigrationSettings.Redo {
		if config.MigrationSettings.Redo {
			step = config.MigrationSettings.Step
		}

		n, err := migrate.ExecMax(sqlPool, config.SqlSettings.DriverName, migrations, migrate.Up, step)
		if err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		log.Printf("Applied %d migrations", n)
	}

	if config.MigrationSettings.Seed {
		if err = ApplySeed(ctx, pool); err != nil {
			return fmt.Errorf("failed to apply database seed: %w", err)
		}
	}

	return nil
}

// ApplySeed initializes the database with data
func ApplySeed(ctx context.Context, pool SqlExecutor) error {
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

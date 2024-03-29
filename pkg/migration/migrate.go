// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package migration

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"path"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	migrate "github.com/rubenv/sql-migrate"
	"go.megpoid.dev/go-skel/pkg/sql"
)

type Options struct {
	TableName      string
	Redo           bool
	Reset          bool
	Rollback       bool
	Step           int
	MigrationAsset AssetOptions
}

type AssetOptions struct {
	FS   embed.FS
	Root string
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, opts Options) error {
	migrations := migrate.EmbedFileSystemMigrationSource{
		FileSystem: opts.MigrationAsset.FS,
		Root:       opts.MigrationAsset.Root,
	}

	migration := migrate.MigrationSet{
		TableName: opts.TableName,
	}

	if opts.Reset {
		_, err := pool.Exec(ctx, "DROP SCHEMA IF EXISTS public CASCADE")
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, "CREATE SCHEMA public")
		if err != nil {
			return nil
		}
		slog.Info("Recreated 'public' schema")
	}

	step := 0

	db := stdlib.OpenDBFromPool(pool)

	if !opts.Reset && (opts.Rollback || opts.Redo) {
		step = opts.Step
		n, err := migration.ExecMaxContext(ctx, db, "postgres", migrations, migrate.Down, step)
		if err != nil {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		slog.Info("Reverted migrations", slog.Int("count", n))
	}

	if opts.Reset || !opts.Rollback || opts.Redo {
		if opts.Redo {
			step = opts.Step
		}

		n, err := migration.ExecMaxContext(ctx, db, "postgres", migrations, migrate.Up, step)
		if err != nil {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		slog.Info("Applied migrations", slog.Int("count", n))
	}

	return nil
}

// ApplySQLFiles initializes the database with data from a directory of SQL files
func ApplySQLFiles(ctx context.Context, conn sql.Executor, config AssetOptions) error {
	assets := config.FS
	entries, err := assets.ReadDir(config.Root)
	if err != nil {
		return err
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func(tx sql.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Failed to rollback transaction", slog.String("error", err.Error()))
		}
	}(tx, ctx)

	for _, entry := range entries {
		if entry.IsDir() || path.Ext(entry.Name()) != ".sql" {
			continue
		}

		seed, err := assets.ReadFile(path.Join(config.Root, entry.Name()))
		if err != nil {
			return err
		}

		slog.Debug("Applying file to database", slog.String("name", entry.Name()))
		if _, err := tx.Exec(ctx, string(seed)); err != nil {
			return fmt.Errorf("failed to apply file %s: %w", entry.Name(), err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

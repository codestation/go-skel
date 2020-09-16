package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"megpoid.xyz/go/go-skel/pkg/sql/helper"
)

// Config hold the database configuration
type Config struct {
	MaxOpenConnections int
}

// NewDatabase opens a database connection indicated by the dsn connection string
func NewDatabase(ctx context.Context, dsn string, config Config) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	db.MapperFunc(helper.ToSnakeCase)
	db.SetMaxOpenConns(config.MaxOpenConnections)

	return db, nil
}

func NewListener(dsn, channel string) (*pq.Listener, error) {
	result := pq.NewListener(dsn, 10*time.Second, time.Minute, func(event pq.ListenerEventType, err error) {
		//_, _ = fmt.Fprintf(os.Stderr, "event in %s listener: %v: %s", channel, event, err)
	})
	err := result.Listen(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to listen to channel: %w", err)
	}
	return result, nil
}

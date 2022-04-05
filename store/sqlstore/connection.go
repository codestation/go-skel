package sqlstore

import (
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"megpoid.xyz/go/go-skel/config"
)

const (
	pingMaxAttempts = 6
	pingTimeoutSecs = 10
)

func (ss *SqlStore) setupConnection(cfg *config.Config) SQLConn {
	db, err := sqlx.Open(cfg.DBAdapter, cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to open database, aborting")
	}

	// total waiting time: 1 minute
	for i := 0; i < pingMaxAttempts; i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutSecs*time.Second)
			defer cancel()

			return db.PingContext(ctx)
		}()

		if err == nil {
			break
		}

		if i < pingMaxAttempts {
			log.Printf("Failed to ping database: %s, retrying in %d seconds", err.Error(), pingTimeoutSecs)
			time.Sleep(pingTimeoutSecs * time.Second)
		} else {
			log.Fatal("Failed to ping database, aborting")
		}
	}

	db.MapperFunc(ToSnakeCase)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return NewDb(db)
}

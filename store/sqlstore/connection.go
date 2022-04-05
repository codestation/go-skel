package sqlstore

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

const (
	pingMaxAttempts = 6
	pingTimeoutSecs = 10
)

const (
	postgresUniqueViolationCode = "23505"
)

func (ss *SqlStore) setupConnection() SQLConn {
	db, err := sqlx.Open(ss.settings.DriverName, ss.settings.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to open database, aborting: %s", err.Error())
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
	db.SetMaxIdleConns(ss.settings.MaxIdleConns)
	db.SetMaxOpenConns(ss.settings.MaxOpenConns)
	db.SetConnMaxLifetime(ss.settings.ConnMaxLifetime)
	db.SetConnMaxIdleTime(ss.settings.ConnMaxIdleTime)

	return NewDb(db)
}

func IsUniqueError(err error, opts ...Option) bool {
	var pqErr *pq.Error
	var sqErr *sqlite3.Error

	switch {
	case errors.As(err, &pqErr):
		if pqErr.Code == postgresUniqueViolationCode {
			for _, opt := range opts {
				if !opt(pqErr) {
					return false
				}
			}
			return true
		}
	case errors.As(err, &sqErr):
		if sqErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			for _, opt := range opts {
				if !opt(pqErr) {
					return false
				}
			}
			return true
		}
	}

	return false
}

type Option func(err error) bool

func WithConstraintName(name string) Option {
	return func(err error) bool {
		var pqErr *pq.Error
		var sqErr *sqlite3.Error

		switch {
		case errors.As(err, &pqErr):
			if pqErr.Constraint == name {
				return true
			}
		case errors.As(err, &sqErr):
			if strings.Contains(sqErr.Error(), name) {
				return true
			}
		}

		return false
	}
}

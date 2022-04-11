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
	"errors"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	pingMaxAttempts = 6
	pingTimeoutSecs = 10
)

const (
	postgresUniqueViolationCode = "23505"
)

func (ss *SqlStore) setupConnection() sqlDb {
	config, err := pgxpool.ParseConfig(ss.settings.DataSourceName)
	if err != nil {
		log.Fatalf("Failed to configure database, aborting: %s", err.Error())
	}

	config.MaxConnLifetime = ss.settings.ConnMaxLifetime
	config.MaxConnIdleTime = ss.settings.ConnMaxIdleTime
	config.MaxConns = int32(ss.settings.MaxOpenConns)
	config.MinConns = int32(ss.settings.MaxIdleConns)

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to open database, aborting: %s", err.Error())
	}

	// total waiting time: 1 minute
	for i := 0; i < pingMaxAttempts; i++ {
		err := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), pingTimeoutSecs*time.Second)
			defer cancel()

			return db.Ping(ctx)
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

	return newPgxWrapper(db)
}

func IsUniqueError(err error, opts ...Option) bool {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		if pgErr.Code == postgresUniqueViolationCode {
			for _, opt := range opts {
				if !opt(pgErr) {
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
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):
			if pgErr.ConstraintName == name {
				return true
			}
		}

		return false
	}
}

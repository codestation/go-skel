// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxListener struct {
	conn *pgxpool.Conn
}

func (l *PgxListener) WaitForNotification(ctx context.Context) (*pgconn.Notification, error) {
	n, err := l.conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (l *PgxListener) Release() {
	l.conn.Release()
}

func NewListener(ctx context.Context, db Acquirer, name string) (*PgxListener, error) {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("LISTEN \"%s\"", name))
	if err != nil {
		return nil, err
	}

	return &PgxListener{conn}, nil
}

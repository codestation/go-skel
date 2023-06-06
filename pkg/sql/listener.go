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

type Notification pgconn.Notification

type Listener struct {
	conn *pgxpool.Conn
}

func (l *Listener) WaitForNotification(ctx context.Context) (*Notification, error) {
	n, err := l.conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return nil, err
	}

	return (*Notification)(n), nil
}

func (l *Listener) Release() {
	l.conn.Release()
}

func NewListener(ctx context.Context, db *PgxWrapper, name string) (*Listener, error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(ctx, fmt.Sprintf("LISTEN \"%s\"", name))
	if err != nil {
		return nil, err
	}

	return &Listener{conn}, nil
}

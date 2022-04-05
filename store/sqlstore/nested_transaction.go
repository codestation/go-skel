/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlstore

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

var uuidGen uuid.Generator

var (
	ErrNoTransactionFound = errors.New("no transaction found")
)

// Driver handles database connections from sqlx.tx or sqlx.tx
type Driver interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
}

// SQLConn adds transaction support for sql databases
type SQLConn interface {
	Driver
	PingContext(ctx context.Context) error

	Close() error
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (SQLConn, error)
	EndTx(txFunc func() error) error
	Commit() error
	Rollback() error
	Tx() *sqlx.Tx
	DB() *sqlx.DB
}

type postgresDatabase struct {
	Driver
	db         *sqlx.DB
	tx         *sqlx.Tx
	savePoints []string
}

func (sd postgresDatabase) PingContext(ctx context.Context) error {
	return sd.db.PingContext(ctx)
}

func NewDb(db *sqlx.DB) SQLConn {
	conn := &postgresDatabase{
		db:     db,
		Driver: db,
	}

	return conn
}

func NewDbFromTx(tx *sqlx.Tx) SQLConn {
	conn := &postgresDatabase{
		tx:     tx,
		Driver: tx,
	}

	return conn
}

func (sd postgresDatabase) Close() error {
	return sd.db.Close()
}

func (sd postgresDatabase) Tx() *sqlx.Tx {
	return sd.tx
}

func (sd postgresDatabase) DB() *sqlx.DB {
	return sd.db
}

// BeginTxx starts a database transaction.
func (sd postgresDatabase) BeginTxx(ctx context.Context, opts *sql.TxOptions) (SQLConn, error) {
	var err error
	if sd.tx == nil {
		sd.tx, err = sd.db.BeginTxx(ctx, opts)
		sd.Driver = sd.tx
	} else {
		var id uuid.UUID
		id, err = uuidGen.NewV4()
		if err != nil {
			return nil, err
		}
		savepoint := "s_" + hex.EncodeToString(id.Bytes())
		sd.savePoints = append(sd.savePoints, savepoint)
		_, err = sd.tx.ExecContext(ctx, "SAVEPOINT "+savepoint)
	}

	if err != nil {
		return nil, err
	}

	return &sd, nil
}

// EndTx wraps a database transaction.
func (sd *postgresDatabase) EndTx(txFunc func() error) (err error) {
	defer func() {
		if p := recover(); p != nil {
			_ = sd.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = sd.Rollback() // err is non-nil; don't change it
		} else {
			err = sd.Commit() // if Commit returns error update err with commit err
		}
	}()
	err = txFunc()
	return
}

// Commit confirms the previous database queries.
func (sd *postgresDatabase) Commit() error {
	if sd.tx == nil {
		return ErrNoTransactionFound
	}
	var err error
	if len(sd.savePoints) > 0 {
		savepoint := sd.savePoints[len(sd.savePoints)-1]
		_, err = sd.tx.Exec("RELEASE SAVEPOINT " + savepoint)
		if err != nil {
			sd.savePoints = sd.savePoints[:len(sd.savePoints)-1]
		}
	} else {
		err = sd.tx.Commit()
	}
	if err != nil {
		return err
	}

	if len(sd.savePoints) == 0 {
		sd.tx = nil
		sd.Driver = nil
	}
	return nil
}

// Rollback cancels the current transaction.
func (sd *postgresDatabase) Rollback() error {
	if sd.tx == nil {
		return nil
	}
	var err error
	if len(sd.savePoints) > 0 {
		savepoint := sd.savePoints[len(sd.savePoints)-1]
		_, err = sd.tx.Exec("ROLLBACK TO SAVEPOINT " + savepoint)
		if err == nil {
			sd.savePoints = sd.savePoints[:len(sd.savePoints)-1]
		}
	} else {
		err = sd.tx.Rollback()
	}
	if err != nil {
		return err
	}

	if len(sd.savePoints) == 0 {
		sd.tx = nil
		sd.Driver = nil
	}
	return nil
}

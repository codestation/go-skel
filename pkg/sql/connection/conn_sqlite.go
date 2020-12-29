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

package connection

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type sqliteDatabase struct {
	Driver
	db      *sqlx.DB
	tx      *sqlx.Tx
	uuidGen uuid.Generator
}

func (sd sqliteDatabase) PingContext(ctx context.Context) error {
	return sd.db.PingContext(ctx)
}

func NewSqlite(db *sqlx.DB) SQLConnection {
	return &sqliteDatabase{db: db, Driver: db, uuidGen: uuid.NewGen()}
}

// TxBegin starts a database transaction.
func (sd sqliteDatabase) TxBegin(ctx context.Context) (SQLConnection, error) {
	var err error
	if sd.tx == nil {
		sd.tx, err = sd.db.BeginTxx(ctx, nil)
		sd.Driver = sd.tx
	} else {
		return nil, errors.New("nested transactions are unsupported")
	}

	if err != nil {
		return nil, err
	}

	return &sd, nil
}

// TxEnd wraps a database transaction.
func (sd *sqliteDatabase) TxEnd(txFunc func() error) (err error) {
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
func (sd *sqliteDatabase) Commit() error {
	if sd.tx == nil {
		return errors.New("no transaction found")
	}
	err := sd.tx.Commit()
	if err != nil {
		return err
	}
	sd.tx = nil
	sd.Driver = nil
	return nil
}

// Rollback cancels the current transaction.
func (sd *sqliteDatabase) Rollback() error {
	if sd.tx == nil {
		return nil
	}
	err := sd.tx.Rollback()
	if err != nil {
		return err
	}
	sd.tx = nil
	sd.Driver = nil
	return nil
}

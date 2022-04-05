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
	"encoding/hex"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresDatabase struct {
	Driver
	db                *sqlx.DB
	tx                *sqlx.Tx
	savepointName     string
	nestedTransaction bool
	uuidGen           uuid.Generator
}

func (sd postgresDatabase) PingContext(ctx context.Context) error {
	return sd.db.PingContext(ctx)
}

func NewPostgresConn(db *sqlx.DB) SQLConnection {
	return &postgresDatabase{db: db, Driver: db, uuidGen: uuid.NewGen()}
}

func (sd postgresDatabase) Close() error {
	return sd.db.Close()
}

// TxBegin starts a database transaction.
func (sd postgresDatabase) TxBegin(ctx context.Context) (SQLConnection, error) {
	var err error
	if sd.tx == nil {
		sd.tx, err = sd.db.BeginTxx(ctx, nil)
		sd.Driver = sd.tx
	} else {
		var id uuid.UUID
		sd.nestedTransaction = true
		id, err = sd.uuidGen.NewV4()
		if err != nil {
			return nil, err
		}
		sd.savepointName = "s_" + hex.EncodeToString(id.Bytes())
		_, err = sd.tx.ExecContext(ctx, "SAVEPOINT "+sd.savepointName)
	}

	if err != nil {
		return nil, err
	}

	return &sd, nil
}

// TxEnd wraps a database transaction.
func (sd *postgresDatabase) TxEnd(txFunc func() error) (err error) {
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
		return errors.New("no transaction found")
	}
	var err error
	if sd.savepointName != "" {
		_, err = sd.tx.Exec("RELEASE SAVEPOINT " + sd.savepointName)
	} else if !sd.nestedTransaction {
		err = sd.tx.Commit()
	}
	if err != nil {
		return err
	}
	sd.tx = nil
	sd.Driver = nil
	return nil
}

// Rollback cancels the current transaction.
func (sd *postgresDatabase) Rollback() error {
	if sd.tx == nil {
		return nil
	}
	var err error
	if sd.savepointName != "" {
		_, err = sd.tx.Exec("ROLLBACK TO SAVEPOINT " + sd.savepointName)
	} else if !sd.nestedTransaction {
		err = sd.tx.Rollback()
	}
	if err != nil {
		return err
	}
	sd.tx = nil
	sd.Driver = nil
	return nil
}

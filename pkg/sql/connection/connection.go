package connection

import (
	"context"
	"encoding/hex"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

// Driver handles database connections from sqlx.tx or sqlx.tx
type Driver interface {
	sqlx.ExecerContext
	sqlx.QueryerContext
	sqlx.PreparerContext
}

// SQLConnection adds transaction support for sql databases
type SQLConnection interface {
	Driver
	TxBegin(ctx context.Context) (SQLConnection, error)
	TxEnd(txFunc func() error) error
	Commit() error
	Rollback() error
	PingContext(ctx context.Context) error
}

type sqlDatabase struct {
	Driver
	db                *sqlx.DB
	tx                *sqlx.Tx
	savepointName     string
	nestedTransaction bool
	uuidGen           uuid.Generator
}

func (sd sqlDatabase) PingContext(ctx context.Context) error {
	return sd.db.PingContext(ctx)
}

func New(db *sqlx.DB) SQLConnection {
	return &sqlDatabase{db: db, Driver: db, uuidGen: uuid.NewGen()}
}

// TxBegin starts a database transaction.
func (sd sqlDatabase) TxBegin(ctx context.Context) (SQLConnection, error) {
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
func (sd *sqlDatabase) TxEnd(txFunc func() error) (err error) {
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
func (sd *sqlDatabase) Commit() error {
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
func (sd *sqlDatabase) Rollback() error {
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

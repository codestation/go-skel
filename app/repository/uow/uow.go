// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package uow

import (
	"context"

	"megpoid.dev/go/go-skel/app/repository"
	"megpoid.dev/go/go-skel/pkg/sql"
)

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name UnitOfWorkStore
type UnitOfWorkStore interface {
	Profiles() repository.ProfileRepo
}

// uowStore has all the repositories of the application
type uowStore struct {
	profiles repository.ProfileRepo
}

func newUowStore(conn sql.Executor) *uowStore {
	return &uowStore{
		profiles: repository.NewProfile(conn),
	}
}

func (u uowStore) Profiles() repository.ProfileRepo {
	return u.profiles
}

type UnitOfWorkBlock func(UnitOfWork) error

//go:generate go run github.com/vektra/mockery/v2@v2.23.1 --name UnitOfWork
type UnitOfWork interface {
	Do(ctx context.Context, fn UnitOfWorkBlock) error
	Store() UnitOfWorkStore
}

type unitOfWork struct {
	conn  sql.Executor
	store *uowStore
}

func New(conn sql.Executor) UnitOfWork {
	return &unitOfWork{
		conn:  conn,
		store: newUowStore(conn),
	}
}

func (u *unitOfWork) Store() UnitOfWorkStore {
	return u.store
}

func (u *unitOfWork) Do(ctx context.Context, fn UnitOfWorkBlock) error {
	err := u.conn.BeginFunc(ctx, func(conn sql.Executor) error {
		uowTx := New(conn)
		if err := fn(uowTx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

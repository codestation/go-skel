// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"megpoid.dev/go/go-skel/pkg/repo"
	"megpoid.dev/go/go-skel/pkg/sql"
)

type HealthcheckRepoImpl struct {
	db sql.Pinger
}

func NewHealthcheckRepo(conn sql.Pinger) *HealthcheckRepoImpl {
	s := &HealthcheckRepoImpl{
		db: conn,
	}

	return s
}

// Execute returns an error if the database doesn't respond
func (s HealthcheckRepoImpl) Execute(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return repo.NewRepoError(repo.ErrBackend, err)
	}
	return nil
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"megpoid.dev/go/go-skel/repository/sqlrepo"
)

type HealthcheckRepoImpl struct {
	db sqlrepo.SqlPinger
}

func NewHealthcheckRepo(conn sqlrepo.SqlPinger) HealthcheckRepo {
	s := &HealthcheckRepoImpl{
		db: conn,
	}

	return s
}

// Execute returns an error if the database doesn't respond
func (s HealthcheckRepoImpl) Execute(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return NewRepoError(ErrBackend, err)
	}
	return nil
}
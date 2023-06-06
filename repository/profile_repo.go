// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/sql"
	"megpoid.dev/go/go-skel/repository/filter"
)

type ProfileRepoImpl struct {
	*GenericStoreImpl[*model.Profile]
}

func NewProfileRepo(conn sql.Executor) ProfileRepo {
	s := &ProfileRepoImpl{
		GenericStoreImpl: NewStore(conn, WithFilters[*model.Profile](
			filter.Rule{
				Key:  "first_name",
				Type: "string",
			}, filter.Rule{
				Key:  "last_name",
				Type: "string",
			}, filter.Rule{
				Key:  "created_at",
				Type: "timestamp",
			},
		)),
	}
	return s
}

func (s *ProfileRepoImpl) GetByEmail(ctx context.Context, email string) (*model.Profile, error) {
	queryBuilder := s.builder.From(s.table).
		Select(s.selectFields...).
		Where(goqu.Ex{"email": email})

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return nil, NewRepoError(ErrBackend, err)
	}

	var result model.Profile
	err = s.db.Get(ctx, &result, query, args...)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, NewRepoError(ErrNotFound, err)
	case err != nil:
		return nil, NewRepoError(ErrBackend, err)
	default:
		return &result, nil
	}
}

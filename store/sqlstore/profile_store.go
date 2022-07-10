// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store"
	"megpoid.xyz/go/go-skel/store/filter"
)

type SqlProfileStore struct {
	*genericStore[*model.Profile]
}

func newSqlProfileStore(sqlStore *SqlStore) store.ProfileStore {
	s := &SqlProfileStore{
		genericStore: NewStore(sqlStore, WithFilterConfig[*model.Profile](filter.Config{Rules: []filter.Rule{{
			Key:  "first_name",
			Type: "string",
		}, {
			Key:  "last_name",
			Type: "string",
		}, {
			Key:  "created_at",
			Type: "timestamp",
		},
		}})),
	}
	return s
}

func (s *SqlProfileStore) GetByUserToken(ctx context.Context, userToken string) (*model.Profile, error) {
	query := s.builder.From(s.table).
		Select(s.selectFields...).
		Where(goqu.Ex{"user_token": userToken})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	var result model.Profile
	err = s.db.Get(ctx, &result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return nil, store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return &result, nil
	}
}

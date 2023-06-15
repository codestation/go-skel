// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/repo"
	"megpoid.dev/go/go-skel/pkg/repo/filter"
	"megpoid.dev/go/go-skel/pkg/sql"
)

type ProfileRepoImpl struct {
	*repo.GenericStoreImpl[*model.Profile]
}

func NewProfile(conn sql.Executor) *ProfileRepoImpl {
	s := &ProfileRepoImpl{
		GenericStoreImpl: repo.NewStore(conn, repo.WithFilters[*model.Profile](
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
	return s.GetBy(ctx, repo.Expr{"email": email})
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"go.megpoid.dev/go-skel/app/model"
	"go.megpoid.dev/go-skel/pkg/repo"
	"go.megpoid.dev/go-skel/pkg/repo/filter"
	"go.megpoid.dev/go-skel/pkg/sql"
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
	return s.GetBy(ctx, repo.Ex{"email": email})
}

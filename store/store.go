// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"context"
	"github.com/gofrs/uuid"
	"megpoid.dev/go/go-skel/model"
	"megpoid.dev/go/go-skel/model/response"
	"megpoid.dev/go/go-skel/store/clause"
)

type Expr map[string]any

// Store lists all the other stores
type Store interface {
	HealthCheck() HealthCheckStore
	Profile() ProfileStore
	WithTransaction(ctx context.Context, f func(s Store) error) error
}

type GenericStore[T model.Modelable] interface {
	Get(ctx context.Context, id model.ID) (T, error)
	GetBy(ctx context.Context, expr Expr) (T, error)
	GetByExternalID(ctx context.Context, externalID uuid.UUID) (T, error)
	List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListBy(ctx context.Context, expr Expr, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListByIds(ctx context.Context, ids []model.ID) (*response.ListResponse[T], error)
	ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error
	ListByEach(ctx context.Context, expr Expr, fn func(item T) error, opts ...clause.FilterOption) error
	Save(ctx context.Context, req T) error
	Update(ctx context.Context, req T) error
	Delete(ctx context.Context, id model.ID) error
	DeleteBy(ctx context.Context, expr Expr) error
	DeleteByExternalId(ctx context.Context, externalId uuid.UUID) error
	Each(ctx context.Context, fn func(entry T) error, opts ...clause.FilterOption) error
}

// HealthCheckStore handles all healthCheck related operations on the store
type HealthCheckStore interface {
	HealthCheck(ctx context.Context) error
}

type ProfileStore interface {
	GenericStore[*model.Profile]
	GetByUserToken(ctx context.Context, userToken string) (*model.Profile, error)
}

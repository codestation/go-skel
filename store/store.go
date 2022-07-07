// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"context"
	"github.com/gofrs/uuid"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/response"
	"megpoid.xyz/go/go-skel/store/clause"
)

// Store lists all the other stores
type Store interface {
	HealthCheck() HealthCheckStore
	Profile() ProfileStore
	WithTransaction(ctx context.Context, f func(s Store) error) error
}

type GenericStore[T any, PT model.Modelable[T]] interface {
	Get(ctx context.Context, id model.ID) (PT, error)
	GetByExternalID(ctx context.Context, externalID uuid.UUID) (PT, error)
	List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T, PT], error)
	ListByIds(ctx context.Context, ids []model.ID) ([]PT, error)
	Save(ctx context.Context, req PT) error
	Update(ctx context.Context, req PT) error
	Delete(ctx context.Context, id model.ID) error
	DeleteByExternalId(ctx context.Context, externalId uuid.UUID) error
}

// HealthCheckStore handles all healthCheck related operations on the store
type HealthCheckStore interface {
	HealthCheck(ctx context.Context) error
}

type ProfileStore interface {
	GenericStore[model.Profile, *model.Profile]
	GetByUserToken(ctx context.Context, userToken string) (*model.Profile, error)
}

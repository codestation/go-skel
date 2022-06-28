// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package store

import (
	"context"
	"github.com/gofrs/uuid"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/response"
)

// Store lists all the other stores
type Store interface {
	HealthCheck() HealthCheckStore
	Profile() ProfileStore
	WithTransaction(ctx context.Context, f func(s Store) error) error
}

type CrudStore[T any, PT model.Modelable[T]] interface {
	Get(ctx context.Context, id model.ID) (PT, error)
	GetByExtID(ctx context.Context, externalID uuid.UUID) (PT, error)
	List(ctx context.Context, opts ...FilterOption) (*response.ListResponse[T], error)
	Save(ctx context.Context, req PT) error
	Update(ctx context.Context, req PT) error
	Delete(ctx context.Context, id model.ID) error
	DeleteByExtId(ctx context.Context, externalId uuid.UUID) error
}

// HealthCheckStore handles all healthCheck related operations on the store
type HealthCheckStore interface {
	HealthCheck(ctx context.Context) error
}

type ProfileStore interface {
	CrudStore[model.Profile, *model.Profile]
	GetByUserToken(ctx context.Context, userToken string) (*model.Profile, error)
}

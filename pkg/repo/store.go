// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repo

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"megpoid.dev/go/go-skel/pkg/clause"
	"megpoid.dev/go/go-skel/pkg/model"
	"megpoid.dev/go/go-skel/pkg/response"
)

type (
	Ex         = goqu.Ex
	ExOr       = goqu.ExOr
	Expression = goqu.Expression
	Op         = goqu.Op
)

var (
	// I represent a schema, table, column or any combination separated by "."
	I = goqu.I
	// C represents a column
	C = goqu.C
	// And represent multiple AND operations toghether
	And = goqu.And
	// Or represent multiple OR operations toghether
	Or = goqu.Or
)

type GenericStore[T model.Modelable] interface {
	Find(ctx context.Context, dest T, id int64) error
	Get(ctx context.Context, id int64) (T, error)
	GetBy(ctx context.Context, expr Expression) (T, error)
	Exists(ctx context.Context, expr Expression) (bool, error)
	List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListBy(ctx context.Context, expr Expression, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListByIds(ctx context.Context, ids []int64) (*response.ListResponse[T], error)
	ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error
	ListByEach(ctx context.Context, expr Expression, fn func(item T) error, opts ...clause.FilterOption) error
	Insert(ctx context.Context, req T) error
	// Upsert inserts a new record in the database, if the target column has a conflict then updates the fields instead
	Upsert(ctx context.Context, req T, target string) (bool, error)
	// Update updates a record on the repository
	Update(ctx context.Context, req T) error
	// UpdateMap updates a record from the repository, only updates the specified fields in the map
	UpdateMap(ctx context.Context, id int64, req map[string]any) error
	// Delete removes a record from the repository, returns ErrNotFound if the ID doesn't exist
	Delete(ctx context.Context, id int64) error
	// DeleteBy removes the records matched by the expression, returns the deleted count on success
	DeleteBy(ctx context.Context, expr Ex) (int64, error)
}

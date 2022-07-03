// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/model/response"
	"megpoid.xyz/go/go-skel/store"
	"megpoid.xyz/go/go-skel/store/filter"
	"megpoid.xyz/go/go-skel/store/paginator"
)

// compile time validator for the interfaces
var (
	_ store.CrudStore[model.Model, *model.Model] = &crudStore[model.Model, *model.Model]{}
)

type crudStore[T any, PT model.Modelable[T]] struct {
	*SqlStore
	table           string
	paginatorConfig paginator.Config
	filterConfig    filter.Config
	selectFields    []any
	defaultFilters  exp.ExpressionList
}

type CrudOption[T any, PT model.Modelable[T]] func(c *crudStore[T, PT])

func WithPaginatorConfig[T any, PT model.Modelable[T]](cfg paginator.Config) CrudOption[T, PT] {
	return func(c *crudStore[T, PT]) {
		c.paginatorConfig = cfg
	}
}

func WithFilterConfig[T any, PT model.Modelable[T]](cfg filter.Config) CrudOption[T, PT] {
	return func(c *crudStore[T, PT]) {
		c.filterConfig = cfg
	}
}

func WithSelectFields[T any, PT model.Modelable[T]](fields ...any) CrudOption[T, PT] {
	return func(c *crudStore[T, PT]) {
		c.selectFields = fields
	}
}

func WithFilters[T any, PT model.Modelable[T]](filters exp.ExpressionList) CrudOption[T, PT] {
	return func(c *crudStore[T, PT]) {
		c.defaultFilters = filters
	}
}

func NewCrudStore[T any, PT model.Modelable[T]](sqlStore *SqlStore, opts ...CrudOption[T, PT]) *crudStore[T, PT] {
	st := &crudStore[T, PT]{SqlStore: sqlStore}
	var defaults []CrudOption[T, PT]
	defaults = append(defaults, WithSelectFields[T, PT]("*"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.table = model.GetTableName[T, PT](new(T))
	return st
}

func (s *crudStore[T, PT]) Get(ctx context.Context, id model.ID) (PT, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"id": id})
	if !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	var result T
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return nil, store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return &result, nil
	}
}

func (s *crudStore[T, PT]) GetByExtID(ctx context.Context, externalID uuid.UUID) (PT, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex{"external_id": externalID})
	if !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	var result T
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return nil, store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return &result, nil
	}
}

func (s *crudStore[T, PT]) List(ctx context.Context, opts ...store.FilterOption) (*response.ListResponse[T], error) {
	query := s.builder.From(s.table).Select(s.selectFields...)
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	p := paginator.New(&s.paginatorConfig)
	f := filter.New(&s.filterConfig)
	for _, opt := range opts {
		opt(p, f)
	}

	query, err := f.Apply(query)
	if err != nil {
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	results := make([]*T, 0)
	cur, err := p.Paginate(ctx, s.db, query, &results)

	switch {
	case errors.Is(err, ErrNoRows):

		return response.NewListResponse(results, cur), nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	default:
		return response.NewListResponse(results, cur), nil
	}
}

func (s *crudStore[T, PT]) Save(ctx context.Context, req PT) error {
	query := s.builder.Insert(s.table).Rows(req).Returning("id")

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrDuplicated, err)
	}

	var id model.ID
	err = s.db.Get(ctx, &id, sql, args...)

	if err != nil {
		if IsUniqueError(err) {
			return store.NewRepoError(store.ErrDuplicated, err)
		}
		return store.NewRepoError(store.ErrBackend, err)
	}

	req.SetID(id)

	return nil
}

func (s *crudStore[T, PT]) Update(ctx context.Context, req PT) error {
	query := s.builder.Update(s.table).Set(req).Where(goqu.Ex{"id": req.GetID()})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrDuplicated, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)

	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}

func (s *crudStore[T, PT]) Delete(ctx context.Context, id model.ID) error {
	query := s.builder.Delete(s.table).Where(goqu.Ex{"id": id})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrDuplicated, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}

func (s *crudStore[T, PT]) DeleteByExtId(ctx context.Context, externalId uuid.UUID) error {
	query := s.builder.Delete(s.table).Where(goqu.Ex{"external_id": externalId})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrDuplicated, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
	}

	if n != 1 {
		return store.NewRepoError(store.ErrNotFound, nil)
	}

	return nil
}

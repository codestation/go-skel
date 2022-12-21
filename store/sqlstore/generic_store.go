// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"megpoid.dev/go/go-skel/model"
	"megpoid.dev/go/go-skel/model/response"
	"megpoid.dev/go/go-skel/store"
	"megpoid.dev/go/go-skel/store/clause"
	"megpoid.dev/go/go-skel/store/filter"
	"megpoid.dev/go/go-skel/store/paginator"
)

// compile time validator for the interfaces
var (
	_ store.GenericStore[*model.Model] = &GenericStore[*model.Model]{}
)

type AttachFunc[T model.Modelable] func(ctx context.Context, results []T, include string) error

type GenericStore[T model.Modelable] struct {
	*SqlStore
	table          string
	selectFields   []any
	defaultFilters exp.ExpressionList
	sortKeys       []string
	includes       []string
	rules          []filter.Rule
	options        []paginator.Option
	attachFunc     AttachFunc[T]
}

type StoreOption[T model.Modelable] func(c *GenericStore[T])

func WithPaginatorOptions[T model.Modelable](opts ...paginator.Option) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.options = opts
	}
}

func WithSelectFields[T model.Modelable](fields ...any) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.selectFields = fields
	}
}

func WithExpressions[T model.Modelable](filters ...exp.Expression) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.defaultFilters = exp.NewExpressionList(exp.AndType, filters...)
	}
}

func WithSortKeys[T model.Modelable](keys ...string) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.sortKeys = keys
	}
}

func WithFilters[T model.Modelable](rules ...filter.Rule) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.rules = rules
	}
}

func WithIncludes[T model.Modelable](includes ...string) StoreOption[T] {
	return func(c *GenericStore[T]) {
		c.includes = includes
	}
}

func NewStore[T model.Modelable](sqlStore *SqlStore, opts ...StoreOption[T]) *GenericStore[T] {
	st := &GenericStore[T]{SqlStore: sqlStore}
	var defaults []StoreOption[T]
	defaults = append(defaults, WithSelectFields[T]("*"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.table = model.GetTableName[T](*new(T))
	return st
}

func (s *GenericStore[T]) zero() T {
	var result T
	return result
}

func (s *GenericStore[T]) new() T {
	return reflect.New(reflect.TypeOf(s.zero()).Elem()).Interface().(T)
}

func (s *GenericStore[T]) AttachFunc(fn AttachFunc[T]) {
	s.attachFunc = fn
}

func (s *GenericStore[T]) Get(ctx context.Context, id model.ID) (T, error) {
	return s.GetBy(ctx, store.Expr{"id": id})
}

func (s *GenericStore[T]) GetBy(ctx context.Context, expr store.Expr) (T, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex(expr))
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	}

	result := s.new()
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, ErrNoRows):
		return s.zero(), store.NewRepoError(store.ErrNotFound, err)
	case err != nil:
		return s.zero(), store.NewRepoError(store.ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *GenericStore[T]) GetByExternalID(ctx context.Context, externalID uuid.UUID) (T, error) {
	return s.GetBy(ctx, store.Expr{"external_id": externalID})
}

func (s *GenericStore[T]) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, store.Expr{}, opts...)
}

func (s *GenericStore[T]) ListBy(ctx context.Context, expr store.Expr, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex(expr))
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	cl := clause.NewClause(
		clause.WithConfig(s.options),
		clause.WithPaginatorKeys(s.sortKeys),
		clause.WithAllowedIncludes(s.includes),
		clause.WithAllowedFilters(s.rules),
	)
	cl.ApplyOptions(opts...)

	results := make([]T, 0)
	cur, err := cl.ApplyFilters(ctx, s.db, query, &results)

	switch {
	case errors.Is(err, ErrNoRows):
		return response.NewListResponse[T](results, cur), nil
	case err != nil:
		return nil, store.NewRepoError(store.ErrBackend, err)
	}

	if s.attachFunc != nil {
		if err := cl.Includes(func(include string) error {
			return s.attachFunc(ctx, results, include)
		}); err != nil {
			return nil, store.NewRepoError(store.ErrBackend, err)
		}
	}

	return response.NewListResponse[T](results, cur), nil
}

func (s *GenericStore[T]) ListByIds(ctx context.Context, ids []model.ID) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, store.Expr{"id": ids})
}

func (s *GenericStore[T]) ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error {
	return s.ListByEach(ctx, store.Expr{}, fn, opts...)
}

func (s *GenericStore[T]) ListByEach(ctx context.Context, expr store.Expr, fn func(item T) error, opts ...clause.FilterOption) error {
	filters := make([]clause.FilterOption, 0, len(opts))
	copy(filters, opts)

	for {
		resp, err := s.ListBy(ctx, expr, filters...)
		if err != nil {
			return err
		}

		for _, e := range resp.Data {
			if err = fn(e); err != nil {
				return err
			}
		}

		if resp.Meta.Next() {
			if resp.Meta.CurrentPage != nil {
				*resp.Meta.CurrentPage += 1
			}

			filters = filters[:0] // clear slice, keep capacity
			filters = append(filters, clause.WithMeta(resp.Meta))
			filters = append(filters, opts...)
		} else {
			return nil
		}
	}
}

func (s *GenericStore[T]) Save(ctx context.Context, req T) error {
	query := s.builder.Insert(s.table).Rows(req).Returning("id")

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
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

func (s *GenericStore[T]) Update(ctx context.Context, req T) error {
	query := s.builder.Update(s.table).Set(req).Where(goqu.Ex{"id": req.GetID()})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
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

func (s *GenericStore[T]) Delete(ctx context.Context, id model.ID) error {
	return s.DeleteBy(ctx, store.Expr{"id": id})
}

func (s *GenericStore[T]) DeleteBy(ctx context.Context, expr store.Expr) error {
	query := s.builder.Delete(s.table).Where(goqu.Ex(expr))

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return store.NewRepoError(store.ErrBackend, err)
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

func (s *GenericStore[T]) DeleteByExternalId(ctx context.Context, externalId uuid.UUID) error {
	return s.DeleteBy(ctx, store.Expr{"external_id": externalId})
}

func (s *GenericStore[T]) Each(ctx context.Context, fn func(entry T) error, opts ...clause.FilterOption) error {
	return s.ListByEach(ctx, store.Expr{}, fn, opts...)
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"
	"errors"
	"reflect"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/clause"
	"megpoid.dev/go/go-skel/pkg/paginator"
	"megpoid.dev/go/go-skel/pkg/response"
	"megpoid.dev/go/go-skel/repository/filter"
	"megpoid.dev/go/go-skel/repository/sqlrepo"
)

// compile time validator for the interfaces
var (
	_ GenericStore[*model.Model] = &GenericStoreImpl[*model.Model]{}
)

type AttachFunc[T model.Modelable] func(ctx context.Context, results []T, include string) error

type GenericStoreImpl[T model.Modelable] struct {
	db             sqlrepo.SqlExecutor
	builder        goqu.DialectWrapper
	prefix         string
	table          string
	selectFields   []any
	defaultFilters exp.ExpressionList
	sortKeys       []string
	includes       []string
	rules          []filter.Rule
	options        []paginator.Option
	attachFunc     AttachFunc[T]
}

type StoreOption[T model.Modelable] func(c *GenericStoreImpl[T])

func WithPaginatorOptions[T model.Modelable](opts ...paginator.Option) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.options = opts
	}
}

func WithSelectFields[T model.Modelable](fields ...any) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.selectFields = fields
	}
}

func WithExpressions[T model.Modelable](filters ...exp.Expression) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.defaultFilters = exp.NewExpressionList(exp.AndType, filters...)
	}
}

func WithSortKeys[T model.Modelable](keys ...string) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.sortKeys = keys
	}
}

func WithFilters[T model.Modelable](rules ...filter.Rule) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.rules = rules
	}
}

func WithIncludes[T model.Modelable](includes ...string) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.includes = includes
	}
}

func WithTablePrefix[T model.Modelable](prefix string) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.prefix = prefix
	}
}

func NewStore[T model.Modelable](conn sqlrepo.SqlExecutor, opts ...StoreOption[T]) *GenericStoreImpl[T] {
	st := &GenericStoreImpl[T]{db: conn}
	st.builder = sqlrepo.NewQueryBuilder()
	var defaults []StoreOption[T]
	defaults = append(defaults, WithSelectFields[T]("*"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.table = st.prefix + model.GetTableName[T](*new(T))
	return st
}

func (s *GenericStoreImpl[T]) zero() T {
	var result T
	return result
}

func (s *GenericStoreImpl[T]) new() T {
	return reflect.New(reflect.TypeOf(s.zero()).Elem()).Interface().(T)
}

func (s *GenericStoreImpl[T]) AttachFunc(fn AttachFunc[T]) {
	s.attachFunc = fn
}

func (s *GenericStoreImpl[T]) Get(ctx context.Context, id model.ID) (T, error) {
	result, err := s.GetBy(ctx, Expr{"id": id})
	if err != nil {
		return s.zero(), err
	}

	if reflect.ValueOf(result).IsNil() {
		return s.zero(), NewRepoError(ErrNotFound, nil)
	}

	return result, nil
}

func (s *GenericStoreImpl[T]) GetBy(ctx context.Context, expr Expr) (T, error) {
	query := s.builder.From(s.table).Select(s.selectFields...).Where(goqu.Ex(expr))
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		query = query.Where(s.defaultFilters)
	}

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), NewRepoError(ErrBackend, err)
	}

	result := s.new()
	err = s.db.Get(ctx, result, sql, args...)

	switch {
	case errors.Is(err, sqlrepo.ErrNoRows):
		return s.zero(), nil
	case err != nil:
		return s.zero(), NewRepoError(ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *GenericStoreImpl[T]) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, Expr{}, opts...)
}

func (s *GenericStoreImpl[T]) ListBy(ctx context.Context, expr Expr, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
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
	case errors.Is(err, sqlrepo.ErrNoRows):
		return response.NewListResponse[T](results, cur), nil
	case err != nil:
		return nil, NewRepoError(ErrBackend, err)
	}

	if s.attachFunc != nil {
		if err := cl.Includes(func(include string) error {
			return s.attachFunc(ctx, results, include)
		}); err != nil {
			return nil, NewRepoError(ErrBackend, err)
		}
	}

	return response.NewListResponse[T](results, cur), nil
}

func (s *GenericStoreImpl[T]) ListByIds(ctx context.Context, ids []model.ID) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, Expr{"id": ids})
}

func (s *GenericStoreImpl[T]) ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error {
	return s.ListByEach(ctx, Expr{}, fn, opts...)
}

func (s *GenericStoreImpl[T]) ListByEach(ctx context.Context, expr Expr, fn func(item T) error, opts ...clause.FilterOption) error {
	filters := make([]clause.FilterOption, 0, len(opts))
	copy(filters, opts)

	for {
		resp, err := s.ListBy(ctx, expr, filters...)
		if err != nil {
			return err
		}

		for _, e := range resp.Items {
			if err = fn(e); err != nil {
				return err
			}
		}

		if resp.Pagination.Next() {
			if resp.Pagination.CurrentPage != nil {
				*resp.Pagination.CurrentPage += 1
			}

			filters = filters[:0] // clear slice, keep capacity
			filters = append(filters, clause.WithMeta(resp.Pagination))
			filters = append(filters, opts...)
		} else {
			return nil
		}
	}
}

func (s *GenericStoreImpl[T]) Save(ctx context.Context, req T) error {
	query := s.builder.Insert(s.table).Rows(req).Returning("id")

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	var id model.ID
	err = s.db.Get(ctx, &id, sql, args...)

	if err != nil {
		if sqlrepo.IsUniqueError(err) {
			return NewRepoError(ErrDuplicated, err)
		}
		return NewRepoError(ErrBackend, err)
	}

	req.SetID(id)

	return nil
}

func (s *GenericStoreImpl[T]) Update(ctx context.Context, req T) error {
	query := s.builder.Update(s.table).Set(req).Where(goqu.Ex{"id": req.GetID()})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)

	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	}

	return nil
}

func (s *GenericStoreImpl[T]) UpdateMap(ctx context.Context, id model.ID, req map[string]any) error {
	query := s.builder.Update(s.table).Set(req).Where(goqu.Ex{"id": id})

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)

	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	}

	return nil
}

func (s *GenericStoreImpl[T]) Upsert(ctx context.Context, req T, target string) (bool, error) {
	conflict := exp.NewDoUpdateConflictExpression(target, req)
	inserted := goqu.Case().When(goqu.L("xmax::text::int").Gt(0), "updated").Else("inserted").As("upsert_status")
	query := s.builder.Insert(s.table).Rows(req).Returning("id", inserted).OnConflict(conflict)

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return false, NewRepoError(ErrBackend, err)
	}

	result := struct {
		ID           model.ID `db:"id"`
		UpsertStatus string   `db:"upsert_status"`
	}{}

	err = s.db.Get(ctx, &result, sql, args...)

	if err != nil {
		if sqlrepo.IsUniqueError(err) {
			return false, NewRepoError(ErrDuplicated, err)
		}
		return false, NewRepoError(ErrBackend, err)
	}

	req.SetID(result.ID)

	return result.UpsertStatus == "inserted", nil
}

func (s *GenericStoreImpl[T]) Delete(ctx context.Context, id model.ID) error {
	if n, err := s.DeleteBy(ctx, Expr{"id": id}); err != nil {
		return err
	} else if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	} else {
		return nil
	}
}

func (s *GenericStoreImpl[T]) DeleteBy(ctx context.Context, expr Expr) (int64, error) {
	query := s.builder.Delete(s.table).Where(goqu.Ex(expr))

	sql, args, err := query.Prepared(true).ToSQL()
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	result, err := s.db.Exec(ctx, sql, args...)
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	return n, nil
}
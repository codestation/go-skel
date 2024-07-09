// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repo

import (
	"context"
	"errors"
	"reflect"
	"slices"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jackc/pgx/v5"
	"github.com/mitchellh/mapstructure"
	"go.megpoid.dev/go-skel/pkg/clause"
	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/repo/filter"
	"go.megpoid.dev/go-skel/pkg/response"
	"go.megpoid.dev/go-skel/pkg/sql"
)

// compile time validator for the interfaces
var (
	_ GenericStore[*model.Model] = &GenericStoreImpl[*model.Model]{}
)

type AttachFunc[T model.Modelable] func(ctx context.Context, results []T, include string) error

type JoinExpression struct {
	Expression exp.Expression
	Condition  exp.JoinCondition
}

type GenericStoreImpl[T model.Modelable] struct {
	Conn           sql.Executor
	Builder        goqu.DialectWrapper
	Table          string
	prefix         string
	selectFields   []any
	joins          []JoinExpression
	returnFields   []any
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

func WithReturnFields[T model.Modelable](fields ...any) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.returnFields = append(c.returnFields, fields...)
	}
}

func WithJoins[T model.Modelable](joins ...JoinExpression) StoreOption[T] {
	return func(c *GenericStoreImpl[T]) {
		c.joins = joins
	}
}

func NewStore[T model.Modelable](conn sql.Executor, opts ...StoreOption[T]) *GenericStoreImpl[T] {
	st := &GenericStoreImpl[T]{Conn: conn}
	st.Builder = sql.NewQueryBuilder()
	var defaults []StoreOption[T]
	defaults = append(defaults, WithSelectFields[T]("*"), WithReturnFields[T]("id"))
	for _, opt := range append(defaults, opts...) {
		opt(st)
	}

	st.Table = st.prefix + model.GetTableName[T](*new(T))
	return st
}

func (s *GenericStoreImpl[T]) WithTx(conn sql.Executor) *GenericStoreImpl[T] {
	return &GenericStoreImpl[T]{
		Conn:           conn,
		Builder:        s.Builder,
		Table:          s.Table,
		prefix:         s.prefix,
		selectFields:   s.selectFields,
		returnFields:   s.returnFields,
		defaultFilters: s.defaultFilters,
		sortKeys:       s.sortKeys,
		includes:       s.includes,
		rules:          s.rules,
		options:        s.options,
		attachFunc:     s.attachFunc,
	}
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

func (s *GenericStoreImpl[T]) First(ctx context.Context, expr Expression, order ...OrderedExpression) (T, error) {
	queryBuilder := s.Builder.From(s.Table).Select(s.selectFields...).Where(expr).
		Order(order...).Limit(1)
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		queryBuilder = queryBuilder.Where(s.defaultFilters)
	}

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), NewRepoError(ErrBackend, err)
	}

	result := s.new()
	err = s.Conn.Get(ctx, result, query, args...)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return s.zero(), NewRepoError(ErrNotFound, nil)
	case err != nil:
		return s.zero(), NewRepoError(ErrBackend, err)
	default:
		return result, nil
	}
}

// Find returns a record from the database. If no record is found then a ErrNotFound is returned
func (s *GenericStoreImpl[T]) Find(ctx context.Context, dest T, id int64) error {
	queryBuilder := s.Builder.From(s.Table).Select(s.selectFields...).Where(goqu.Ex{"id": id})
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		queryBuilder = queryBuilder.Where(s.defaultFilters)
	}

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	err = s.Conn.Get(ctx, dest, query, args...)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return NewRepoError(ErrNotFound, nil)
	case err != nil:
		return NewRepoError(ErrBackend, err)
	default:
		return nil
	}
}

func (s *GenericStoreImpl[T]) CountBy(ctx context.Context, expr Expression) (int64, error) {
	queryBuilder := s.Builder.From(s.Table).Select(goqu.COUNT("*")).Where(expr)
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		queryBuilder = queryBuilder.Where(s.defaultFilters)
	}

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	var count int64
	err = s.Conn.Get(ctx, &count, query, args...)
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	return count, nil
}

// Get returns a record from the database. If no record is found then a ErrNotFound is returned
func (s *GenericStoreImpl[T]) Get(ctx context.Context, id int64) (T, error) {
	return s.GetBy(ctx, Ex{"id": id})
}

// GetBy returns the first record from the database matching the expression.
// If no record is found then a ErrNotFound is returned
func (s *GenericStoreImpl[T]) GetBy(ctx context.Context, expr Expression) (T, error) {
	queryBuilder := s.Builder.From(s.Table).Select(s.selectFields...).Where(expr)
	if s.defaultFilters != nil && !s.defaultFilters.IsEmpty() {
		queryBuilder = queryBuilder.Where(s.defaultFilters)
	}

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), NewRepoError(ErrBackend, err)
	}

	result := s.new()
	err = s.Conn.Get(ctx, result, query, args...)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return s.zero(), NewRepoError(ErrNotFound, nil)
	case err != nil:
		return s.zero(), NewRepoError(ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *GenericStoreImpl[T]) GetForUpdate(ctx context.Context, expr Expression, order ...OrderedExpression) (T, error) {
	queryBuilder := s.Builder.From(s.Table).Select(s.selectFields...).Where(expr)

	queryBuilder = queryBuilder.Order(order...).Limit(1).ForUpdate(goqu.SkipLocked)

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return s.zero(), NewRepoError(ErrBackend, err)
	}

	result := s.new()
	err = s.Conn.Get(ctx, result, query, args...)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return s.zero(), NewRepoError(ErrNotFound, nil)
	case err != nil:
		return s.zero(), NewRepoError(ErrBackend, err)
	default:
		return result, nil
	}
}

func (s *GenericStoreImpl[T]) Exists(ctx context.Context, expr Expression) (bool, error) {
	_, err := s.GetBy(ctx, expr)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *GenericStoreImpl[T]) List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, Ex{}, opts...)
}

func (s *GenericStoreImpl[T]) ListBy(ctx context.Context, expr Expression, opts ...clause.FilterOption) (*response.ListResponse[T], error) {
	query := s.Builder.From(s.Table).Select(s.selectFields...).Where(expr)
	if s.defaultFilters != nil {
		query = query.Where(s.defaultFilters)
	}

	for _, join := range s.joins {
		query = query.Join(join.Expression, join.Condition)
	}

	cl := clause.NewClause(
		clause.WithConfig(s.options),
		clause.WithPaginatorKeys(s.sortKeys),
		clause.WithAllowedIncludes(s.includes),
		clause.WithAllowedFilters(s.rules),
	)
	cl.ApplyOptions(opts...)

	results := make([]T, 0)
	cur, err := cl.ApplyFilters(ctx, s.Conn, query, &results)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
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

func (s *GenericStoreImpl[T]) ListByIDs(ctx context.Context, ids []int64) (*response.ListResponse[T], error) {
	return s.ListBy(ctx, Ex{"id": ids})
}

func (s *GenericStoreImpl[T]) ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error {
	return s.ListByEach(ctx, Ex{}, fn, opts...)
}

func (s *GenericStoreImpl[T]) ListByEach(ctx context.Context, expr Expression, fn func(item T) error, opts ...clause.FilterOption) error {
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

func (s *GenericStoreImpl[T]) Insert(ctx context.Context, req T) error {
	queryBuilder := s.Builder.Insert(s.Table).Rows(req).Returning(s.returnFields...)

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	err = s.Conn.Get(ctx, req, query, args...)
	if err != nil {
		if sql.IsUniqueError(err) {
			return NewRepoError(ErrDuplicated, err)
		}
		return NewRepoError(ErrBackend, err)
	}

	return nil
}

func (s *GenericStoreImpl[T]) Update(ctx context.Context, req T) error {
	queryBuilder := s.Builder.Update(s.Table).Set(req).Where(goqu.Ex{"id": req.GetID()})

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	result, err := s.Conn.Exec(ctx, query, args...)
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	n := result.RowsAffected()

	if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	}

	return nil
}

func (s *GenericStoreImpl[T]) UpdateMap(ctx context.Context, id int64, req map[string]any) error {
	queryBuilder := s.Builder.Update(s.Table).Set(req).Where(goqu.Ex{"id": id})

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	result, err := s.Conn.Exec(ctx, query, args...)
	if err != nil {
		return NewRepoError(ErrBackend, err)
	}

	n := result.RowsAffected()

	if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	}

	return nil
}

func (s *GenericStoreImpl[T]) UpdateMapBy(ctx context.Context, req map[string]any, expr Expression) (int64, error) {
	queryBuilder := s.Builder.Update(s.Table).Set(req).Where(expr)

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	result, err := s.Conn.Exec(ctx, query, args...)
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	n := result.RowsAffected()

	return n, nil
}

func (s *GenericStoreImpl[T]) Upsert(ctx context.Context, req T, target string) (bool, error) {
	var conflict exp.ConflictExpression
	if target != "" {
		conflict = exp.NewDoUpdateConflictExpression(target, req)
	} else {
		conflict = exp.NewDoNothingConflictExpression()
	}

	inserted := goqu.Case().When(goqu.L("xmax::text::int").Gt(0), "updated").Else("inserted").As("upsert_status")
	queryBuilder := s.Builder.Insert(s.Table).Rows(req).Returning(append(s.returnFields, inserted)...).OnConflict(conflict)

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return false, NewRepoError(ErrBackend, err)
	}

	result := map[string]any{}
	err = s.Conn.Get(ctx, &result, query, args...)
	if err != nil {
		if sql.IsUniqueError(err) {
			return false, NewRepoError(ErrDuplicated, err)
		}
		return false, NewRepoError(ErrBackend, err)
	}

	// do not use mapstructure if there is only one field and it's the id
	if slices.Equal(s.returnFields, []any{"id"}) {
		switch value := result["id"].(type) {
		case int32:
			req.SetID(int64(value))
		case int64:
			req.SetID(value)
		default:
			return false, NewRepoError(ErrBackend, err)
		}
	} else {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Squash: true, Result: req})
		if err != nil {
			return false, NewRepoError(ErrBackend, err)
		}

		if err := decoder.Decode(result); err != nil {
			return false, NewRepoError(ErrBackend, err)
		}
	}

	return result["upsert_status"] == "inserted", nil
}

func (s *GenericStoreImpl[T]) Delete(ctx context.Context, id int64) error {
	if n, err := s.DeleteBy(ctx, Ex{"id": id}); err != nil {
		return err
	} else if n != 1 {
		return NewRepoError(ErrNotFound, nil)
	} else {
		return nil
	}
}

func (s *GenericStoreImpl[T]) DeleteBy(ctx context.Context, expr Ex) (int64, error) {
	queryBuilder := s.Builder.Delete(s.Table).Where(expr)

	query, args, err := queryBuilder.Prepared(true).ToSQL()
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	result, err := s.Conn.Exec(ctx, query, args...)
	if err != nil {
		return 0, NewRepoError(ErrBackend, err)
	}

	n := result.RowsAffected()

	return n, nil
}

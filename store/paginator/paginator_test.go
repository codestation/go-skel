// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import (
	"context"
	"errors"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

type User struct {
	ID   int32
	Name string
}

func TestPaginatorDefault(t *testing.T) {
	p := New()
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	query, err := p.PaginateDataset(query, &User{})
	if assert.NoError(t, err) {
		sql, args, err := query.Prepared(true).ToSQL()
		if assert.NoError(t, err) {
			assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "id" ASC LIMIT $1`, sql)
			assert.Equal(t, []any{int64(100 + 1)}, args)
		}
	}
}

func TestPaginatorMultiple(t *testing.T) {
	p := New(&Config{
		Keys: []string{"Name", "ID"},
	})
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	query, err := p.PaginateDataset(query, &User{})
	if assert.NoError(t, err) {
		sql, args, err := query.Prepared(true).ToSQL()
		if assert.NoError(t, err) {
			assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "id" ASC LIMIT $1`, sql)
			assert.Equal(t, []any{int64(100 + 1)}, args)
		}
	}
}

func TestPaginatorMultipleRule(t *testing.T) {
	p := New(&Config{
		Rules: []Rule{
			{Key: "Name"},
			{Key: "ID", SQLRepr: "users.id"},
		},
	})
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	query, err := p.PaginateDataset(query, &User{})
	if assert.NoError(t, err) {
		sql, args, err := query.Prepared(true).ToSQL()
		if assert.NoError(t, err) {
			assert.Equal(t, `SELECT "id", "name" FROM "users" ORDER BY "name" ASC, "users"."id" ASC LIMIT $1`, sql)
			assert.Equal(t, []any{int64(100 + 1)}, args)
		}
	}
}

func TestPaginatorCursor(t *testing.T) {
	users := []*User{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
		{ID: 3, Name: "C"},
	}

	p1 := New(&Config{
		Keys:  []string{"Name", "ID"},
		Limit: 2,
	})

	cursor1, err := p1.PaginateResults(&users)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(users))
		if assert.NotNil(t, cursor1.After) {
			assert.Equal(t, "WyJCIiwyXQ==", *cursor1.After)
		}
	}

	users2 := []*User{
		{ID: 3, Name: "C"},
	}

	p2 := New(&Config{
		Keys:  []string{"Name", "ID"},
		Limit: 2,
		After: *cursor1.After,
	})

	cursor2, err := p2.PaginateResults(&users2)
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(users2))
		assert.Nil(t, cursor2.After)
		if assert.NotNil(t, cursor2.Before) {
			assert.Equal(t, "WyJDIiwzXQ==", *cursor2.Before)
		}
	}
}

func TestPaginatorMultipleCursor(t *testing.T) {
	p := New(&Config{
		Keys:  []string{"Name", "ID"},
		Limit: 2,
		After: "WyJCIiwyXQ==",
	})
	query := goqu.Dialect("postgres").From("users").Select("id", "name")
	query, err := p.PaginateDataset(query, &User{})
	if assert.NoError(t, err) {
		sql, args, err := query.Prepared(true).ToSQL()
		if assert.NoError(t, err) {
			assert.Equal(t, `SELECT "id", "name" FROM "users" WHERE (("name" > $1) OR (("name" = $2) AND ("id" > $3))) ORDER BY "name" ASC, "id" ASC LIMIT $4`, sql)
			assert.Equal(t, []any{"B", "B", int64(2), int64(3)}, args)
		}
	}
}

type sqlSelector struct {
	Error error
	Users []*User
}

func (s *sqlSelector) Select(_ context.Context, dest any, _ string, _ ...any) error {
	if s.Error != nil {
		return s.Error
	}
	pDest := dest.(*[]*User)
	*pDest = s.Users
	return nil
}

func TestPaginatorPaginate(t *testing.T) {
	p := New(&Config{
		Keys:  []string{"Name", "ID"},
		Limit: 2,
	})

	query := goqu.Dialect("postgres").From("users").Select("id", "name")

	results := make([]*User, 0)
	selector := &sqlSelector{
		Users: []*User{
			{ID: 1, Name: "A"},
			{ID: 2, Name: "B"},
			{ID: 3, Name: "C"},
		},
	}
	c, err := p.Paginate(context.Background(), selector, query, &results)
	if assert.NoError(t, err) {
		assert.Nil(t, c.Before)
		if assert.NotNil(t, c.After) {
			assert.Equal(t, "WyJCIiwyXQ==", *c.After)
		}
	}
}

func TestPaginatorPaginateError(t *testing.T) {
	p := New(&Config{})

	query := goqu.Dialect("postgres").From("users")
	paginatorErr := errors.New("an error ocurred")

	results := make([]*User, 0)
	selector := &sqlSelector{
		Error: paginatorErr,
	}

	_, err := p.Paginate(context.Background(), selector, query, &results)
	assert.ErrorIs(t, err, paginatorErr)
}

func TestPaginatorPaginateDatasetError(t *testing.T) {
	p := New(&Config{})

	query := goqu.Dialect("postgres").From("users")
	var results struct{}

	_, err := p.PaginateDataset(query, &results)
	assert.ErrorIs(t, err, ErrInvalidModel)
}

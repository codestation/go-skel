package repo

import (
	"context"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/clause"
	"megpoid.dev/go/go-skel/pkg/response"
)

type Expr map[string]any

type GenericStore[T model.Modelable] interface {
	Get(ctx context.Context, id model.ID) (T, error)
	GetBy(ctx context.Context, expr Expr) (T, error)
	List(ctx context.Context, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListBy(ctx context.Context, expr Expr, opts ...clause.FilterOption) (*response.ListResponse[T], error)
	ListByIds(ctx context.Context, ids []model.ID) (*response.ListResponse[T], error)
	ListEach(ctx context.Context, fn func(item T) error, opts ...clause.FilterOption) error
	ListByEach(ctx context.Context, expr Expr, fn func(item T) error, opts ...clause.FilterOption) error
	Save(ctx context.Context, req T) error
	// Upsert inserts a new record in the database, if the target column has a conflict then updates the fields instead
	Upsert(ctx context.Context, req T, target string) (bool, error)
	// Update updates a record on the repository
	Update(ctx context.Context, req T) error
	// UpdateMap updates a record from the repository, only updates the specified fields in the map
	UpdateMap(ctx context.Context, id model.ID, req map[string]any) error
	// Delete removes a record from the repository, returns ErrNotFound if the ID doesn't exist
	Delete(ctx context.Context, id model.ID) error
	// DeleteBy removes the records matched by the expression, returns the deleted count on success
	DeleteBy(ctx context.Context, expr Expr) (int64, error)
}

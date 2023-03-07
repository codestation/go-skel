// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"math"

	"megpoid.dev/go/go-skel/app/model"
	"megpoid.dev/go/go-skel/pkg/paginator"
)

type ListResponse[T model.Modelable] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination,omitempty"`
}

func NewListResponse[T model.Modelable](results []T, c *paginator.Cursor) *ListResponse[T] {
	switch c.Type() {
	case paginator.MetaCursor:
		cur := c.Cursor()
		return &ListResponse[T]{
			Items: results,
			Pagination: Pagination{
				Type: string(c.Type()),
				PaginationCursor: PaginationCursor{
					NextCursor: cur.After,
					PrevCursor: cur.Before,
				},
			},
		}
	case paginator.MetaOffset:
		off := c.Offset()
		return &ListResponse[T]{
			Items: results,
			Pagination: Pagination{
				Type: string(c.Type()),
				PaginationOffset: PaginationOffset{
					TotalRecords:   model.NewType(off.Total),
					CurrentPage:    model.NewType(off.Page),
					MaxPage:        model.NewType(int(math.Ceil(float64(off.Total) / float64(off.ItemsPerPage)))),
					RecordsPerPage: model.NewType(off.ItemsPerPage),
				},
			},
		}
	default:
		return &ListResponse[T]{
			Items: results,
		}
	}
}

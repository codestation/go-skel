// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"math"

	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/types"
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
					TotalRecords:   types.AsPointer(off.Total),
					CurrentPage:    types.AsPointer(off.Page),
					MaxPage:        types.AsPointer(int(math.Ceil(float64(off.Total) / float64(off.ItemsPerPage)))),
					RecordsPerPage: types.AsPointer(off.ItemsPerPage),
				},
			},
		}
	default:
		return &ListResponse[T]{
			Items: results,
		}
	}
}

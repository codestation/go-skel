// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"megpoid.dev/go/go-skel/model"
	"megpoid.dev/go/go-skel/store/paginator"
)

type ListResponse[T model.Modelable] struct {
	Data []T `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func NewListResponse[T model.Modelable](results []T, c *paginator.Cursor) *ListResponse[T] {
	switch c.Type() {
	case paginator.MetaCursor:
		cur := c.Cursor()
		return &ListResponse[T]{
			Data: results,
			Meta: CursorPagination{
				Items:      len(results),
				NextCursor: cur.After,
				PrevCursor: cur.Before,
			},
		}
	case paginator.MetaOffset:
		off := c.Offset()
		return &ListResponse[T]{
			Data: results,
			Meta: OffsetPagination{
				Items: len(results),
				Total: off.Total,
				Page:  off.Page,
			},
		}
	default:
		return &ListResponse[T]{
			Data: results,
		}
	}
}

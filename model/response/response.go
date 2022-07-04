// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store/paginator/cursor"
)

type ListResponse[T any, PT model.Modelable[T]] struct {
	Data []PT       `json:"data"`
	Meta Pagination `json:"meta"`
}

func NewListResponse[T any, PT model.Modelable[T]](results []PT, c *cursor.Cursor) *ListResponse[T, PT] {
	return &ListResponse[T, PT]{
		Data: results,
		Meta: Pagination{
			Items:      len(results),
			NextCursor: c.After,
			PrevCursor: c.Before,
		},
	}
}

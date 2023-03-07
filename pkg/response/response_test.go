// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/app/model"
	paginator2 "megpoid.dev/go/go-skel/pkg/paginator"
	"megpoid.dev/go/go-skel/pkg/paginator/cursor"
)

func TestNewListResponseEmpty(t *testing.T) {
	results := make([]*model.Profile, 0)
	response := NewListResponse(results, &paginator2.Cursor{})
	assert.Equal(t, 0, len(response.Items))
	assert.NotNil(t, response.Pagination)
}

func TestNewListResponseCursor(t *testing.T) {
	results := []*model.Profile{
		model.NewProfile(),
	}

	cur := &paginator2.Cursor{}
	cur.SetCursor(&cursor.Cursor{
		After:  model.NewType("after"),
		Before: model.NewType("before"),
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Items))
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, "cursor", response.Pagination.Type)
	assert.Equal(t, "after", *response.Pagination.PaginationCursor.NextCursor)
	assert.Equal(t, "before", *response.Pagination.PaginationCursor.PrevCursor)
}

func TestNewListResponseOffset(t *testing.T) {
	results := []*model.Profile{
		model.NewProfile(),
	}

	cur := &paginator2.Cursor{}
	cur.SetOffset(&paginator2.Page{
		Items:        1,
		Total:        10,
		Page:         4,
		ItemsPerPage: 3,
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Items))
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, "offset", response.Pagination.Type)
	assert.Equal(t, 4, *response.Pagination.PaginationOffset.MaxPage)
	assert.Equal(t, 10, *response.Pagination.PaginationOffset.TotalRecords)
	assert.Equal(t, 4, *response.Pagination.PaginationOffset.CurrentPage)
	assert.Equal(t, 3, *response.Pagination.PaginationOffset.RecordsPerPage)
}

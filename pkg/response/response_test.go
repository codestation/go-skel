// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.megpoid.dev/go-skel/pkg/model"
	"go.megpoid.dev/go-skel/pkg/paginator"
	"go.megpoid.dev/go-skel/pkg/paginator/cursor"
	"go.megpoid.dev/go-skel/pkg/types"
)

type Profile struct {
	model.Model
}

func TestNewListResponseEmpty(t *testing.T) {
	results := make([]*Profile, 0)
	response := NewListResponse(results, &paginator.Cursor{})
	assert.Equal(t, 0, len(response.Items))
	assert.NotNil(t, response.Pagination)
}

func TestNewListResponseCursor(t *testing.T) {
	results := []*Profile{{}}

	cur := &paginator.Cursor{}
	cur.SetCursor(&cursor.Cursor{
		After:  types.AsPointer("after"),
		Before: types.AsPointer("before"),
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Items))
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, "cursor", response.Pagination.Type)
	assert.Equal(t, "after", *response.Pagination.NextCursor)
	assert.Equal(t, "before", *response.Pagination.PrevCursor)
}

func TestNewListResponseOffset(t *testing.T) {
	results := []*Profile{{}}

	cur := &paginator.Cursor{}
	cur.SetOffset(&paginator.Page{
		Items:        1,
		Total:        10,
		Page:         4,
		ItemsPerPage: 3,
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Items))
	assert.NotNil(t, response.Pagination)
	assert.Equal(t, "offset", response.Pagination.Type)
	assert.Equal(t, 4, *response.Pagination.MaxPage)
	assert.Equal(t, 10, *response.Pagination.TotalRecords)
	assert.Equal(t, 4, *response.Pagination.CurrentPage)
	assert.Equal(t, 3, *response.Pagination.RecordsPerPage)
}

// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.megpoid.dev/go-skel/pkg/types"
)

func TestMetaEmpty(t *testing.T) {
	meta := Pagination{}
	assert.False(t, meta.Next())
}

func TestMetaCursor(t *testing.T) {
	meta := Pagination{
		Type: "cursor",
		PaginationCursor: PaginationCursor{
			NextCursor: types.AsPointer("cursor"),
		},
	}

	assert.True(t, meta.Next())

	data, err := meta.MarshalJSON()
	assert.NoError(t, err)

	expected := map[string]any{
		"type":        "cursor",
		"next_cursor": "cursor",
		"prev_cursor": nil,
	}

	expectedJSON, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJSON), string(data))
}

func TestMetaPage(t *testing.T) {
	meta := Pagination{
		Type: "page",
		PaginationOffset: PaginationOffset{
			CurrentPage:    types.AsPointer(1),
			MaxPage:        types.AsPointer(2),
			TotalRecords:   types.AsPointer(2),
			RecordsPerPage: types.AsPointer(1),
		},
	}

	assert.True(t, meta.Next())

	data, err := meta.MarshalJSON()
	assert.NoError(t, err)

	expected := map[string]any{
		"type":             "page",
		"current_page":     1,
		"max_page":         2,
		"records_per_page": 1,
		"total":            2,
	}

	expectedJSON, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJSON), string(data))
}

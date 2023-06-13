// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/pkg/model"
)

func TestMetaEmpty(t *testing.T) {
	meta := Pagination{}
	assert.False(t, meta.Next())
}

func TestMetaCursor(t *testing.T) {
	meta := Pagination{
		Type: "cursor",
		PaginationCursor: PaginationCursor{
			NextCursor: model.NewType("cursor"),
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

	expectedJson, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJson), string(data))
}

func TestMetaPage(t *testing.T) {
	meta := Pagination{
		Type: "page",
		PaginationOffset: PaginationOffset{
			CurrentPage:    model.NewType(1),
			MaxPage:        model.NewType(2),
			TotalRecords:   model.NewType(2),
			RecordsPerPage: model.NewType(1),
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

	expectedJson, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJson), string(data))
}

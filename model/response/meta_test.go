package response

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/model"
	"testing"
)

func TestMetaEmpty(t *testing.T) {
	meta := Meta{}
	assert.False(t, meta.Next())
}

func TestMetaCursor(t *testing.T) {
	meta := Meta{
		Items: 1,
		Type:  "cursor",
		CursorMeta: CursorMeta{
			NextCursor: model.NewType("cursor"),
		},
	}

	assert.True(t, meta.Next())

	data, err := meta.MarshalJSON()
	assert.NoError(t, err)

	expected := map[string]any{
		"items":       1,
		"next_cursor": "cursor",
		"prev_cursor": nil,
	}

	expectedJson, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJson), string(data))
}

func TestMetaPage(t *testing.T) {
	meta := Meta{
		Items: 1,
		Type:  "page",
		OffsetMeta: OffsetMeta{
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
		"current_page":     1,
		"items":            1,
		"max_page":         2,
		"records_per_page": 1,
		"total":            2,
	}

	expectedJson, _ := json.Marshal(expected)
	assert.JSONEq(t, string(expectedJson), string(data))
}

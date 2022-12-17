package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/go-skel/model"
	"megpoid.dev/go/go-skel/store/paginator"
	"megpoid.dev/go/go-skel/store/paginator/cursor"
)

func TestNewListResponseEmpty(t *testing.T) {
	results := make([]*model.Profile, 0)
	response := NewListResponse(results, &paginator.Cursor{})
	assert.Equal(t, 0, len(response.Data))
	assert.NotNil(t, response.Meta)
}

func TestNewListResponseCursor(t *testing.T) {
	results := []*model.Profile{
		model.NewProfile(),
	}

	cur := &paginator.Cursor{}
	cur.SetCursor(&cursor.Cursor{
		After:  model.NewType("after"),
		Before: model.NewType("before"),
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Data))
	assert.NotNil(t, response.Meta)
	assert.Equal(t, "cursor", response.Meta.Type)
	assert.Equal(t, "after", *response.Meta.CursorMeta.NextCursor)
	assert.Equal(t, "before", *response.Meta.CursorMeta.PrevCursor)
}

func TestNewListResponseOffset(t *testing.T) {
	results := []*model.Profile{
		model.NewProfile(),
	}

	cur := &paginator.Cursor{}
	cur.SetOffset(&paginator.Page{
		Items:        1,
		Total:        10,
		Page:         4,
		ItemsPerPage: 3,
	})

	response := NewListResponse(results, cur)
	assert.Equal(t, 1, len(response.Data))
	assert.NotNil(t, response.Meta)
	assert.Equal(t, "offset", response.Meta.Type)
	assert.Equal(t, 4, *response.Meta.OffsetMeta.MaxPage)
	assert.Equal(t, 10, *response.Meta.OffsetMeta.TotalRecords)
	assert.Equal(t, 4, *response.Meta.OffsetMeta.CurrentPage)
	assert.Equal(t, 3, *response.Meta.OffsetMeta.RecordsPerPage)
}

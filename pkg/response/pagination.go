// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import "encoding/json"

type PaginationCursor struct {
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

type PaginationOffset struct {
	CurrentPage    *int `json:"current_page,omitempty"`
	MaxPage        *int `json:"max_page,omitempty"`
	TotalRecords   *int `json:"total,omitempty"`
	RecordsPerPage *int `json:"records_per_page,omitempty"`
}

type Pagination struct {
	Type string `json:"type"`
	PaginationCursor
	PaginationOffset
}

func (m *Pagination) Next() bool {
	if m.NextCursor != nil {
		return true
	}
	if m.CurrentPage != nil && m.MaxPage != nil && *m.CurrentPage < *m.MaxPage {
		return true
	}

	return false
}

func (m *Pagination) MarshalJSON() ([]byte, error) {
	if m.CurrentPage != nil {
		return json.Marshal(&struct {
			Type string `json:"type"`
			PaginationOffset
		}{
			Type:             m.Type,
			PaginationOffset: m.PaginationOffset,
		})
	} else {
		return json.Marshal(&struct {
			Type string `json:"type"`
			PaginationCursor
		}{
			Type:             m.Type,
			PaginationCursor: m.PaginationCursor,
		})
	}
}

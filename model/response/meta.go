// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import "encoding/json"

type CursorMeta struct {
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

type OffsetMeta struct {
	CurrentPage    *int `json:"current_page,omitempty"`
	MaxPage        *int `json:"max_page,omitempty"`
	TotalRecords   *int `json:"total,omitempty"`
	RecordsPerPage *int `json:"records_per_page,omitempty"`
}

type Meta struct {
	Items int    `json:"items"`
	Type  string `json:"type"`
	CursorMeta
	OffsetMeta
}

func (m *Meta) Next() bool {
	if m.NextCursor != nil {
		return true
	}
	if m.CurrentPage != nil && m.MaxPage != nil && *m.CurrentPage < *m.MaxPage {
		return true
	}

	return false
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	if m.CurrentPage != nil {
		return json.Marshal(&struct {
			Items int `json:"items"`
			//Type  string `json:"type"`
			OffsetMeta
		}{
			Items: m.Items,
			//Type:       m.Type,
			OffsetMeta: m.OffsetMeta,
		})
	} else {
		return json.Marshal(&struct {
			Items int `json:"items"`
			//Type  string `json:"type"`
			CursorMeta
		}{
			Items: m.Items,
			//Type:       m.Type,
			CursorMeta: m.CursorMeta,
		})
	}
}

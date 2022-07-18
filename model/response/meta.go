// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

type CursorPagination struct {
	Items      int     `json:"items"`
	NextCursor *string `json:"next_cursor"`
	PrevCursor *string `json:"prev_cursor"`
}

type OffsetPagination struct {
	Items int   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
}

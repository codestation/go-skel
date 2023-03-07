// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package cursor

// Cursor cursor data
type Cursor struct {
	After  *string `json:"after" query:"after"`
	Before *string `json:"before" query:"before"`
}

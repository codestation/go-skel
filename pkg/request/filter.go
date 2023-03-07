// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package request

type Pagination struct {
	Limit  *int    `query:"limit"`
	After  *string `query:"after"`
	Before *string `query:"before"`
	Page   *int    `query:"page"`
}

type Filter struct {
	Field     string
	Operation string
	Value     string
}

type TypeSort string

const (
	TypeSortAsc  TypeSort = "ASC"
	TypeSortDesc TypeSort = "DESC"
)

type SortEntry struct {
	Field     string
	Direction TypeSort
}

type QueryParams struct {
	Pagination Pagination
	Filters    []Filter
	Includes   []string
	Sort       []SortEntry
	Search     string
}

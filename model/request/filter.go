// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package request

import "strings"

type FilterType string

const (
	FilterEqual          FilterType = "eq"
	FilterNotEqual       FilterType = "neq"
	FilterLessThan       FilterType = "lt"
	FilterLessOrEqual    FilterType = "lte"
	FilterGreaterThan    FilterType = "gt"
	FilterGreaterOrEqual FilterType = "gte"
	FilterHas            FilterType = "has"
	FilterIn             FilterType = "in"
)

type Filter struct {
	Field     string
	Operation FilterType
	Value     string
}

func (f Filter) Values() []string {
	return strings.Split(f.Value, ",")
}

type Pagination struct {
	Limit  *int    `query:"limit"`
	After  *string `query:"after"`
	Before *string `query:"before"`
}

type QueryParams struct {
	Pagination Pagination
	Filters    []Filter
	Search     string
}

// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package request

import "strings"

type OperationType string

const (
	OperationEqual          OperationType = "eq"
	OperationNotEqual       OperationType = "neq"
	OperationLessThan       OperationType = "lt"
	OperationLessOrEqual    OperationType = "lte"
	OperationGreaterThan    OperationType = "gt"
	OperationGreaterOrEqual OperationType = "gte"
	OperationHas            OperationType = "has"
	OperationIn             OperationType = "in"
)

type Filter struct {
	Field     string
	Operation OperationType
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

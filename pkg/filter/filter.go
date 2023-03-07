// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package filter

import (
	"errors"
	"fmt"
	"strings"

	"megpoid.dev/go/go-skel/oapi"
	"megpoid.dev/go/go-skel/pkg/request"
)

type FilterParams struct {
	Before   *oapi.Before
	After    *oapi.After
	Page     *oapi.Page
	Q        *oapi.Query
	Limit    *oapi.Limit
	Includes *oapi.Includes
	Filters  *oapi.Filters
	Fields   *oapi.Fields
	Sort     *oapi.Sort
}

func NewFilterFromParams(params FilterParams) (*request.QueryParams, error) {
	query := &request.QueryParams{
		Pagination: request.Pagination{
			Limit:  params.Limit,
			After:  params.After,
			Before: params.Before,
			Page:   params.Page,
		},
	}

	if params.Includes != nil {
		query.Includes = *params.Includes
	}

	if params.Sort != nil {
		sortArgs := strings.Split(*params.Sort, ",")
		for _, s := range sortArgs {
			if len(s) > 0 {
				switch s[0] {
				case '-':
					query.Sort = append(query.Sort, request.SortEntry{
						Field:     s[1:],
						Direction: request.TypeSortDesc,
					})
				case '+':
					query.Sort = append(query.Sort, request.SortEntry{
						Field:     s[1:],
						Direction: request.TypeSortAsc,
					})
				default:
					query.Sort = append(query.Sort, request.SortEntry{
						Field:     s,
						Direction: request.TypeSortAsc,
					})
				}
			}
		}
	}

	if params.Q != nil {
		query.Search = *params.Q
	}

	var errorList []error

	if params.Filters != nil {
		for key, value := range *params.Filters {
			filterParts := strings.Split(key, "__")
			if len(filterParts) == 2 {
				query.Filters = append(query.Filters, request.Filter{
					Field:     filterParts[0],
					Operation: filterParts[1],
					Value:     value,
				})
			} else {
				errorList = append(errorList, fmt.Errorf("invalid query param: %s", key))
				continue
			}
		}
	}

	if len(errorList) > 0 {
		return nil, errors.Join(errorList...)
	}

	return query, nil
}

// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"megpoid.dev/go/go-skel/model/request"
)

func NewFilter(c echo.Context) (*request.QueryParams, error) {
	query := &request.QueryParams{}
	if err := c.Bind(&query.Pagination); err != nil {
		return nil, fmt.Errorf("failed to bind pagination filter: %w", err)
	}

	var errorList []error
	for key, value := range c.QueryParams() {
		switch key {
		case "after":
			// managed by bind
			fallthrough
		case "before":
			// managed by bind
			fallthrough
		case "page":
			// managed by bind
			fallthrough
		case "limit":
			// managed by bind
		case "includes":
			if len(value) > 0 {
				query.Includes = strings.Split(value[0], ",")
			}
		case "q":
			if len(value) > 0 {
				query.Search = value[0]
			}
		default:
			filterParts := strings.Split(key, "__")
			if len(filterParts) == 2 {
				query.Filters = append(query.Filters, request.Filter{
					Field:     filterParts[0],
					Operation: filterParts[1],
					Value:     value[0], //ignore other repeated filters
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

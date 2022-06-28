// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"megpoid.xyz/go/go-skel/model/request"
	"strings"
)

func NewFilter(c echo.Context) (*request.QueryParams, error) {
	query := &request.QueryParams{}
	if err := c.Bind(&query.Pagination); err != nil {
		return nil, fmt.Errorf("failed to bind pagination filter: %w", err)
	}

	for key, value := range c.QueryParams() {
		switch key {
		case "after":
			fallthrough
		case "before":
			fallthrough
		case "limit":
		// managed by bind
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
				// skip values that aren't filters
				continue
			}
		}
	}

	return query, nil
}

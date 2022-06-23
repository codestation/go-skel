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
	if err := c.Bind(query); err != nil {
		return nil, fmt.Errorf("failed to bind pagination filter: %w", err)
	}

	for key, value := range c.QueryParams() {
		if len(value) == 0 {
			// skip empty filter
			continue
		}
		switch key {
		case "after":
			fallthrough
		case "before":
			fallthrough
		case "limit":
		// managed by bind
		case "q":
			query.Search = value[0]
		default:
			filterParts := strings.Split(key, "__")
			if len(filterParts) == 2 {
				switch request.FilterType(filterParts[1]) {
				case request.FilterEqual:
					fallthrough
				case request.FilterNotEqual:
					fallthrough
				case request.FilterGreaterThan:
					fallthrough
				case request.FilterGreaterOrEqual:
					fallthrough
				case request.FilterLessThan:
					fallthrough
				case request.FilterLessOrEqual:
					fallthrough
				case request.FilterHas:
					fallthrough
				case request.FilterIn:
					query.Filters = append(query.Filters, request.Filter{
						Field:     filterParts[0],
						Operation: request.FilterType(filterParts[1]),
						Value:     value[0], //ignore other repeated filters
					})
				default:
					// ignore filters outside the available ones
					continue
				}
			} else {
				// ignore other fields not separated by __
				continue
			}
		}
	}

	return query, nil
}

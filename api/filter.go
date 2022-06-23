// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/labstack/echo/v4"
	"log"
	"megpoid.xyz/go/go-skel/model/request"
)

func NewFilter(c echo.Context) *request.Filter {
	query := &request.Filter{}
	if err := c.Bind(query); err != nil {
		log.Printf("NewFilter failed: %s", err.Error())
		return &request.Filter{}
	}
	return query
}

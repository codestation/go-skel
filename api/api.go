// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/labstack/echo/v4"
	"megpoid.xyz/go/go-skel/app"
)

const appName = "goapp"
const apiVersion = "v1"

type API struct {
	app  app.IApp
	root *echo.Group
}

func Init(srv *app.Server) (*API, error) {
	api := &API{
		app: app.New(srv),
	}

	api.root = api.app.Srv().EchoServer.Group("/apis/" + appName + "/" + apiVersion)

	// initialize all handlers
	api.InitStatus()
	api.InitProfile()

	return api, nil
}

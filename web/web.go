package web

import (
	"github.com/labstack/echo/v4"
	"megpoid.xyz/go/go-skel/app"
)

type Web struct {
	app  app.IApp
	root *echo.Group
}

func New(srv *app.Server) *Web {
	web := &Web{
		app: app.New(srv),
	}

	web.root = web.app.Srv().EchoServer.Group("")

	// initialize all handlers
	web.InitStatic()

	return web
}

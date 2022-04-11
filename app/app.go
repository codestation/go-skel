package app

import (
	"megpoid.xyz/go/go-skel/model"
	"time"
)

type App struct {
	svr     *Server
	timeNow func() time.Time
}

func New(svr *Server) *App {
	app := &App{svr: svr, timeNow: time.Now}
	return app
}

func (a *App) Srv() *Server {
	return a.svr
}

func (a *App) Config() *model.Config {
	return &a.svr.cfg
}

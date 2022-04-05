package app

import "megpoid.xyz/go/go-skel/model"

type App struct {
	svr *Server
}

func New(svr *Server) *App {
	app := &App{svr: svr}
	return app
}

func (a *App) Srv() *Server {
	return a.svr
}

func (a *App) Config() *model.Config {
	return &a.svr.cfg
}

package app

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

package services

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type httpServer struct {
	HTTPServer *http.Server
	e          *echo.Echo
}

func NewHTTPServer(addr string, e *echo.Echo) *httpServer {
	svr := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return &httpServer{HTTPServer: svr, e: e}
}

func (h *httpServer) Serve() error {
	if err := h.e.StartServer(h.HTTPServer); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}

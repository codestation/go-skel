package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"megpoid.xyz/go/go-skel/model"
	"megpoid.xyz/go/go-skel/store"
	"megpoid.xyz/go/go-skel/store/sqlstore"
)

type Server struct {
	cfg        model.Config
	sqlStore   *sqlstore.SqlStore
	Store      store.Store
	Server     *http.Server
	EchoServer *echo.Echo
}

func NewServer(cfg model.Config) (*Server, error) {
	s := &Server{cfg: cfg}

	// Store initialization, could use a different database or non-sql store
	s.sqlStore = sqlstore.New(cfg.SqlSettings)
	s.Store = s.sqlStore

	// HTTP server initialization
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.GeneralSettings.Debug
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1M"))

	if e.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	}

	s.EchoServer = e

	return s, nil
}

const (
	shutdownTimeout = 30 * time.Second
)

func (s *Server) Start() error {
	s.Server = &http.Server{
		Addr:         s.cfg.ServerSettings.ListenAddress,
		ReadTimeout:  s.cfg.ServerSettings.ReadTimeout,
		WriteTimeout: s.cfg.ServerSettings.WriteTimeout,
		IdleTimeout:  s.cfg.ServerSettings.IdleTimeout,
	}

	log.Printf("Starting server")

	go func() {
		err := s.EchoServer.StartServer(s.Server)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %s", err.Error())
			time.Sleep(time.Second)
		}
	}()

	return nil
}

func (s *Server) StopHTTPServer() {
	if s.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.EchoServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server: StopHTTPServer: Shutdown failed: %s", err.Error())
		}

		if err := s.Server.Close(); err != nil {
			log.Printf("Server: StopHTTPServer: Close failed: %s", err.Error())
		}
		s.Server = nil
	}
}

func (s *Server) Shutdown() {
	s.StopHTTPServer()
	if err := s.Store.Close(); err != nil {
		log.Printf("Server: Store: Close failed: %s", err.Error())
	}
}

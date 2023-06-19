// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v4emb"
	"megpoid.dev/go/go-skel/app/controller"
	"megpoid.dev/go/go-skel/app/repository"
	"megpoid.dev/go/go-skel/app/repository/uow"
	"megpoid.dev/go/go-skel/app/usecase"
	"megpoid.dev/go/go-skel/config"
	"megpoid.dev/go/go-skel/oapi"
	"megpoid.dev/go/go-skel/pkg/i18n"
	"megpoid.dev/go/go-skel/pkg/jwt"
	"megpoid.dev/go/go-skel/pkg/sql"
	"megpoid.dev/go/go-skel/web"
)

const (
	shutdownTimeout = 30 * time.Second
)

type App struct {
	cfg        *config.Config
	conn       sql.Database
	Server     *http.Server
	EchoServer *echo.Echo
}

func NewApp(cfg *config.Config) (*App, error) {
	s := &App{cfg: cfg}

	// Database initialization
	pool, err := sql.NewConnection(sql.Config(cfg.DatabaseSettings))
	if err != nil {
		return nil, err
	}

	s.conn = sql.NewPgxWrapper(pool)

	// Repository initialization (not attached to the unit of work)
	healthcheckRepo := repository.NewHealthCheck(s.conn)

	// Unit of Work initialization (all repos are initialized here)
	unitOfWork := uow.New(s.conn)

	// Usecase initialization
	authUsecase := usecase.NewAuth(cfg.ServerSettings.JwtSecret)
	healthcheckUsecase := usecase.NewHealthcheck(healthcheckRepo)
	profileUsecase := usecase.NewProfile(unitOfWork)

	// Controller initialization
	ctrl := controller.Controller{
		AuthController:        controller.NewAuth(cfg, authUsecase),
		ProfileController:     controller.NewProfile(cfg, profileUsecase),
		HealthcheckController: controller.NewHealthCheck(cfg, healthcheckUsecase),
	}

	// HTTP server initialization
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.GeneralSettings.Debug
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.ServerSettings.CorsAllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/swagger")
		},
	}))
	e.Use(middleware.Recover())
	e.Use(i18n.LoadMessagePrinter("user_lang"))
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(cfg.ServerSettings.BodyLimit))
	e.Use(middleware.RequestID())
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = ErrorHandler(e)
	s.EchoServer = e

	// Serve Swagger UI
	handler := v4emb.NewHandlerWithConfig(swgui.Config{
		Title:       "Skel API",
		SwaggerJSON: "/swagger/docs/openapi.yaml",
		BasePath:    "/swagger",
		SettingsUI: map[string]string{
			"defaultModelsExpandDepth": "1",
		},
	})

	spec, err := oapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("error loading spec: %w", err)
	}

	oapiMiddleware := jwt.OapiValidator(spec, cfg.ServerSettings.JwtSecret)
	e.Use(oapiMiddleware)

	swagger := echo.WrapHandler(handler)
	e.GET("/swagger", swagger)
	e.GET("/swagger/*", swagger)

	// Embed openapi docs
	assetHandler := http.FileServer(http.FS(oapi.Assets()))
	e.GET("/swagger/docs/*", echo.WrapHandler(http.StripPrefix("/swagger/docs/", assetHandler)))

	web.New(e)

	oapi.RegisterHandlersWithBaseURL(e, &ctrl, controller.BaseURL())

	return s, nil
}

func (s *App) Start() error {
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

func (s *App) stopHTTPServer() {
	if s.Server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.EchoServer.Shutdown(shutdownCtx); err != nil {
			log.Printf("App: stopHTTPServer: Shutdown failed: %s", err.Error())
		}

		if err := s.Server.Close(); err != nil {
			log.Printf("App: stopHTTPServer: Close failed: %s", err.Error())
		}

		s.Server = nil
	}
}

func (s *App) Shutdown() {
	s.stopHTTPServer()
	s.conn.Close()
}

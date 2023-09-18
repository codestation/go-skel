// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hibiken/asynq"
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
	"megpoid.dev/go/go-skel/pkg/apperror"
	"megpoid.dev/go/go-skel/pkg/i18n"
	"megpoid.dev/go/go-skel/pkg/jwt"
	"megpoid.dev/go/go-skel/pkg/sql"
	"megpoid.dev/go/go-skel/pkg/task"
	"megpoid.dev/go/go-skel/pkg/validator"
	"megpoid.dev/go/go-skel/web"
)

const (
	shutdownTimeout = 30 * time.Second
)

type Config struct {
	General  config.GeneralSettings
	Database config.DatabaseSettings
	Server   config.ServerSettings
}

type App struct {
	cfg        Config
	conn       sql.Database
	Server     *http.Server
	EchoServer *echo.Echo
}

func NewApp(cfg Config) (*App, error) {
	s := &App{cfg: cfg}

	if cfg.General.Debug {
		textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		slog.SetDefault(slog.New(textHandler))
	}

	// Database initialization
	pool, err := sql.NewConnection(sql.Config(cfg.Database))
	if err != nil {
		return nil, err
	}

	s.conn = sql.NewPgxWrapper(pool)

	// Repository initialization (not attached to the unit of work)
	healthcheckRepo := repository.NewHealthCheck(s.conn)

	// Unit of Work initialization (all repos are initialized here)
	unitOfWork := uow.New(s.conn)

	// Redis client config
	redisClient := asynq.RedisClientOpt{
		Addr: cfg.General.RedisAddr,
	}

	// Usecase initialization
	authUsecase := usecase.NewAuth(cfg.Server.JwtSecret)
	healthcheckUsecase := usecase.NewHealthcheck(healthcheckRepo)
	profileUsecase := usecase.NewProfile(unitOfWork)
	taskUsecase := task.NewClient(redisClient)

	// Controller initialization
	ctrl := controller.Controller{
		AuthController:        controller.NewAuth(cfg.Server, authUsecase),
		ProfileController:     controller.NewProfile(cfg.Server, profileUsecase),
		HealthcheckController: controller.NewHealthCheck(cfg.Server, healthcheckUsecase),
		TaskController:        controller.NewTask(cfg.Server, taskUsecase),
		DelayController:       controller.NewDelay(cfg.Server, taskUsecase),
	}

	// HTTP server initialization
	e := echo.New()
	e.HideBanner = true
	e.Debug = cfg.General.Debug
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.Server.CorsAllowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(ctx echo.Context) bool {
			return strings.HasPrefix(ctx.Path(), controller.BaseURL()+"/swagger")
		},
	}))
	e.Use(middleware.Recover())
	e.Use(i18n.LoadMessagePrinter("user_lang"))
	e.Use(middleware.Logger())
	e.Use(middleware.BodyLimit(cfg.Server.BodyLimit))
	e.Use(middleware.RequestID())
	e.Validator = validator.NewCustomValidator()
	e.HTTPErrorHandler = apperror.ErrorHandler(e)
	s.EchoServer = e

	// Serve Swagger UI
	handler := v4emb.NewHandlerWithConfig(swgui.Config{
		Title:       "Skel API",
		SwaggerJSON: controller.BaseURL() + "/swagger/docs/openapi.yaml",
		BasePath:    controller.BaseURL() + "/swagger",
		SettingsUI: map[string]string{
			"defaultModelsExpandDepth": "1",
		},
	})

	spec, err := oapi.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("error loading spec: %w", err)
	}

	skipperFunc := jwt.WithSkipperFunc(func(ctx echo.Context) bool {
		path := ctx.Path()
		return strings.HasPrefix(path, controller.BaseURL()+"/swagger")
	})

	oapiMiddleware := jwt.OapiValidator(spec, cfg.Server.JwtSecret, skipperFunc)
	e.Use(oapiMiddleware)

	group := e.Group(controller.BaseURL())
	swagger := echo.WrapHandler(handler)
	group.GET("/swagger", swagger)
	group.GET("/swagger/*", swagger)

	// Embed openapi docs
	assetHandler := http.FileServer(http.FS(oapi.Assets()))
	group.GET("/swagger/docs/*", echo.WrapHandler(http.StripPrefix(controller.BaseURL()+"/swagger/docs/", assetHandler)))

	web.New(e)

	oapi.RegisterHandlersWithBaseURL(e, &ctrl, controller.BaseURL())

	return s, nil
}

func (s *App) Start() error {
	s.Server = &http.Server{
		Addr:         s.cfg.Server.ListenAddress,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	slog.Info("Starting server")

	go func() {
		err := s.EchoServer.StartServer(s.Server)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Error starting server", slog.String("error", err.Error()))
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
			slog.Error("App: stopHTTPServer: shutdown failed", slog.String("error", err.Error()))
		}

		if err := s.Server.Close(); err != nil {
			slog.Error("App: stopHTTPServer: close failed", slog.String("error", err.Error()))
		}

		s.Server = nil
	}
}

func (s *App) Shutdown() {
	s.stopHTTPServer()
	s.conn.Close()
}

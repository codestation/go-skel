//go:generate mockgen -destination=mocks/echo.go -package=mocks github.com/labstack/echo/v4 Context
/*
Copyright © 2020 codestation <codestation404@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package delivery

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/domain/usecase"
)

const apiVersion = "v1"

type HTTPHandler struct {
	uc  usecase.UseCase
	cfg *config.Config
}

func (h *HTTPHandler) getClaim(c echo.Context, key string) (string, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	value, ok := claims[key].(string)
	if !ok {
		return "", fmt.Errorf("failed to find value for key %s", key)
	}

	return value, nil
}

func NewHTTPHandler(uc usecase.UseCase, cfg *config.Config) *HTTPHandler {
	return &HTTPHandler{uc: uc, cfg: cfg}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	app := e.Group("/apis/goapp/" + apiVersion)
	app.GET("/status/livez", h.LiveCheck)
	app.GET("/status/readyz", h.ReadyCheck)

	protected := app.Group("")
	protected.Use(middleware.JWT(h.cfg.JWTSecret))
	protected.GET("/hello", h.Hello)
}

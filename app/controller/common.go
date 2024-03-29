// Copyright 2023 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.megpoid.dev/go-skel/config"
	"go.megpoid.dev/go-skel/oapi"
	"go.megpoid.dev/go-skel/pkg/i18n"
	"golang.org/x/text/message"
)

const (
	appName    = "goapp"
	apiVersion = "v1"
)

func BaseURL() string {
	return "/apis/" + appName + "/" + apiVersion
}

var _ oapi.ServerInterface = &Controller{}

type Controller struct {
	AuthController
	ProfileController
	HealthcheckController
	TaskController
	DelayController
}

type common struct {
	config config.ServerSettings
}

func newCommon(config config.ServerSettings) common {
	return common{
		config: config,
	}
}

// printer returns a printer to localize messages to other languages.
func (a *common) printer(ctx echo.Context) *message.Printer {
	return message.NewPrinter(i18n.GetLanguageTags(ctx))
}

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (a *common) UseJWT(g *echo.Group) {
	jwtConfig := echojwt.Config{
		SigningKey: a.config.JwtSecret,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &JwtCustomClaims{}
		},
	}

	g.Use(echojwt.WithConfig(jwtConfig))
}

func (a *common) GetUserID(c echo.Context) string {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	return claims.UserID
}

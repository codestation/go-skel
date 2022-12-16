// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package api

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func (api *API) UseJWT(g *echo.Group) {
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: api.app.Config().GeneralSettings.JwtSecret,
	}

	g.Use(middleware.JWTWithConfig(config))
}

func (api *API) GetClaim(c echo.Context, key string) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims[key].(string)
}

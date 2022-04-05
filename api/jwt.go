package api

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Schema string `json:"schema"`
	Domain string `json:"domain"`
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

func (api *API) getClaim(c echo.Context, key string) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims[key].(string)
}
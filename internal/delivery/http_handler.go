//go:generate mockgen -destination=mocks/echo.go -package=mocks github.com/labstack/echo/v4 Context

package delivery

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"megpoid.xyz/go/go-skel/internal/config"
	"megpoid.xyz/go/go-skel/internal/usecases"
)

type HTTPHandler struct {
	uc  usecase.APIUsecase
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

func NewHTTPHandler(uc usecase.APIUsecase, cfg *config.Config) *HTTPHandler {
	return &HTTPHandler{uc: uc, cfg: cfg}
}

func (h *HTTPHandler) Register(e *echo.Echo) {
	apiGroup := e.Group("/api")
	v1 := apiGroup.Group("/v1")
	v1.GET("/stat/health", h.HealthCheck)

	app := v1.Group("/app")
	app.Use(middleware.JWT(h.cfg.JWTSecret))
}

func (h *HTTPHandler) HealthCheck(c echo.Context) error {
	if err := h.uc.HealthCheck(c.Request().Context()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

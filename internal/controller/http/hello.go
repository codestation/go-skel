package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) Hello(c echo.Context) error {
	accountId, err := h.getClaim(c, "account_id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":    "hello world",
		"account_id": accountId,
	})
}

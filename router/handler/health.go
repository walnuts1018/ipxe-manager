package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Liveness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func (h *Handler) Readiness(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

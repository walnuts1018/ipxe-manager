package handler

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetIpxe(c echo.Context) error {
	script, err := h.usecase.GetIpxeScript(c.Request().Context())
	if err != nil {
		slog.Error("failed to generate iPXE script", slog.String("error", err.Error()))
		return echo.NewHTTPError(500, "failed to generate iPXE script")
	}
	return c.Stream(200, "text/plain", script)
}

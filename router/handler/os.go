package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/ipxe-manager/definitions"
)

func (h *Handler) GetOSStatus(c echo.Context) error {
	os, err := h.usecase.GetOS(c.Request().Context())
	if err != nil {
		slog.Error("failed to get OS status", slog.String("error", err.Error()))
		return echo.NewHTTPError(500, "failed to get OS status")
	}
	return c.JSON(http.StatusOK, map[string]string{"os": os})
}

func (h *Handler) SetOS(c echo.Context) error {
	var params struct {
		OS string `json:"os" form:"os"`
	}
	if err := c.Bind(&params); err != nil {
		slog.Error("failed to bind request", slog.String("error", err.Error()))
		return echo.NewHTTPError(400, "invalid request")
	}

	if params.OS == "" {
		return echo.NewHTTPError(400, "os is required")
	}

	if err := h.usecase.SetOS(c.Request().Context(), params.OS); err != nil {
		if errors.Is(err, definitions.ErrIPXEScriptNotFound) {
			return echo.NewHTTPError(404, "OS script not found")
		}
		slog.Error("failed to set OS", slog.String("error", err.Error()))
		return echo.NewHTTPError(500, "failed to set OS")
	}
	return c.String(http.StatusOK, "ok")
}

package handler

import (
	"errors"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/walnuts1018/ipxe-manager/definitions"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := getAccessToken(c.Request().Header.Get("Authorization"))

		if err := h.usecase.CheckAuthorization(c.Request().Context(), token); err != nil {
			slog.Error("authorization failed", slog.String("error", err.Error()), slog.String("token", token))
			if errors.Is(err, definitions.ErrInvalidToken) {
				return echo.NewHTTPError(401, "invalid token")
			}
			return echo.NewHTTPError(500, "internal server error")
		}

		return next(c)
	}
}

func getAccessToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}

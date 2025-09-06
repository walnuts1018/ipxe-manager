package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/walnuts1018/ipxe-manager/config"
	"github.com/walnuts1018/ipxe-manager/router/handler"
	"github.com/walnuts1018/ipxe-manager/tracer"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func NewRouter(cfg *config.AppConfig, handler *handler.Handler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: skipper,
	}))
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(tracer.ServiceName, otelecho.WithSkipper(
		skipper,
	)))

	e.GET("/livez", handler.Liveness).Name = "health.liveness"
	e.GET("/readyz", handler.Readiness).Name = "health.readiness"

	api := e.Group("/api")
	v1 := api.Group("/v1")
	v1.Use(handler.AuthMiddleware)

	os := v1.Group("/os")
	os.GET("/status", handler.GetOSStatus).Name = "GetOSStatus"
	os.POST("/set", handler.SetOS).Name = "SetOS"

	e.GET("/boot.ipxe", handler.GetIpxe).Name = "GetIpxe"

	return e
}

func skipper(c echo.Context) bool {
	// Skip tracing for health check endpoints
	return c.Path() == "/livez" || c.Path() == "/readyz"
}

package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func Apply(e *echo.Echo, otelTracer trace.Tracer, corsAllowedOrigins string) {
	e.HideBanner = true
	e.HidePort = true

	e.Use(echomw.Recover())
	e.Use(echomw.RequestID())
	e.Use(echomw.Logger())
	e.Use(echomw.Secure())
	e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: parseOrigins(corsAllowedOrigins),
	}))
	e.Use(echomw.RateLimiter(echomw.NewRateLimiterMemoryStore(20)))
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			propagators := otel.GetTextMapPropagator()
			ctx := propagators.Extract(c.Request().Context(), propagation.HeaderCarrier(c.Request().Header))
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	})
}

func parseOrigins(csv string) []string {
	fields := strings.Split(csv, ",")
	origins := make([]string, 0, len(fields))
	for _, f := range fields {
		if o := strings.TrimSpace(f); o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}

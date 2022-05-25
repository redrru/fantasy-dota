package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/log"
)

func LoggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if c.Path() == "/metrics" {
					return
				}
				log.GetLogger().Info(c.Request().Context(),
					"Handle request",
					zap.String("method", c.Request().Method),
					zap.String("path", c.Path()),
					zap.Int("status", c.Response().Status),
					zap.Error(err),
				)
			}()

			err = next(c)
			return
		}
	}
}

package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/redrru/fantasy-dota/pkg/log"
)

func RecoveringMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				r := recover()
				if r == nil {
					return
				}

				if e, ok := r.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("%v", r)
				}

				log.GetLogger().Error(c.Request().Context(), "Panic recovered", zap.String("url", c.Path()), zap.Error(err))
			}()

			err = next(c)
			return
		}
	}
}

package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/redrru/fantasy-dota/pkg/server"
)

// GetExample - Example GET handler.
// (GET /example)
func (s *Server) GetExample(ctx echo.Context) error {
	panic(fmt.Errorf("example panic"))
	//  return echo.NewHTTPError(http.StatusInternalServerError, "example error")
}

// PostExample - Example POST handler.
// (POST /example)
func (s *Server) PostExample(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, server.ExampleResponse{Name: "qwe"})
}

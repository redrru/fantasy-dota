package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
	"github.com/redrru/fantasy-dota/pkg/server"
)

// GetExample - Example GET handler.
// (GET /example)
func (s *Server) GetExample(ctx echo.Context) error {
	models, err := s.usecase.ExampleGet(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "example error")
	}

	var result server.ExampleResponse
	for _, models := range models {
		result = append(result, server.ExampleObject{Name: models.Name})
	}

	return ctx.JSON(http.StatusOK, result)
}

// PostExample - Example POST handler.
// (POST /example)
func (s *Server) PostExample(ctx echo.Context) error {
	req := new(server.PostExampleJSONRequestBody)
	if err := ctx.Bind(req); err != nil {
		return err
	}
	if req.Name == "" {
		return fmt.Errorf("empty name")
	}

	model := entity.ExampleModel{
		Name: req.Name,
	}

	err := s.usecase.ExamplePost(ctx.Request().Context(), model)
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

package main

import (
	"github.com/labstack/echo/v4"

	application "github.com/redrru/fantasy-dota/internal/fantasy-dota"
	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
	"github.com/redrru/fantasy-dota/internal/fantasy-dota/repository"
	"github.com/redrru/fantasy-dota/internal/fantasy-dota/usecase"
	"github.com/redrru/fantasy-dota/internal/gateways/http"
	"github.com/redrru/fantasy-dota/pkg/server"
)

func main() {
	app := application.NewApplication()

	repo := repository.NewRepository(app.DB)
	uc := usecase.NewUsecase(repo)

	e := echo.New()
	server.RegisterHandlers(e, http.NewServer(uc))
	app.RegisterHTTP(e)

	app.RegisterMigrationModel(entity.ExampleModel{})
	// app.RegisterFetchers(fetchers.NewExample())

	app.Run()
}

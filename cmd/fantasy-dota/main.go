package main

import (
	"github.com/labstack/echo/v4"

	application "github.com/redrru/fantasy-dota/internal/fantasy-dota"
	"github.com/redrru/fantasy-dota/internal/fantasy-dota/fetchers"
	"github.com/redrru/fantasy-dota/internal/gateways/http"
	"github.com/redrru/fantasy-dota/pkg/server"
)

func main() {
	app := application.NewApplication()

	e := echo.New()
	server.RegisterHandlers(e, http.NewServer())
	app.RegisterHTTP(e)

	app.RegisterFetchers(fetchers.NewExample())

	app.Run()
}

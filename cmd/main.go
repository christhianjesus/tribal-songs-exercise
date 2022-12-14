package main

import (
	"christhianguevara/songs-search-exercise/cmd/config"
	"christhianguevara/songs-search-exercise/domain/constants"
	"christhianguevara/songs-search-exercise/internal/handlers"
	"christhianguevara/songs-search-exercise/internal/middlewares"
	"christhianguevara/songs-search-exercise/internal/resources"
	"christhianguevara/songs-search-exercise/internal/services"
	"fmt"
	"net/http"

	"github.com/joeshaw/envdecode"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := config.Config{}

	if err := envdecode.Decode(&conf); err != nil {
		panic(fmt.Errorf("cannot read from env: %w", err))
	}

	e := echo.New()
	e.Use(middlewares.KeyAuth(conf.AuthKey))

	setupHandlers(e.Group(constants.PrefixPath))

	e.Logger.Fatal(e.Start(conf.Addr()))
}

func setupHandlers(router *echo.Group) {
	// Client
	client := http.DefaultClient

	// Resources
	iTunesResource := resources.NewITunesResource(client)

	// Services
	songsService := services.NewSongsService([]resources.SongsResource{iTunesResource})

	// Handlers
	songsHandler := handlers.NewSongsHandler(songsService)

	// Register routes
	songsHandler.RegisterRoutes(router)
}

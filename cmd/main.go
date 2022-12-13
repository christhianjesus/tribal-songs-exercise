package main

import (
	"christhianguevara/songs-search-exercise/cmd/config"
	"christhianguevara/songs-search-exercise/domain/constants"
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/handlers"
	"christhianguevara/songs-search-exercise/internal/middlewares"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"fmt"

	"github.com/joeshaw/envdecode"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
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
	// Mock
	songsServiceMock := new(mocks.SongsService)
	songsServiceMock.
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{{ID: "1"}}, nil)

	// Handlers
	songsHandler := handlers.NewSongsHandler(songsServiceMock)

	// Register routes
	songsHandler.RegisterRoutes(router)
}

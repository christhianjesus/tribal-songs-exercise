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

	"github.com/go-redis/redis/v9"
	"github.com/joeshaw/envdecode"
	"github.com/labstack/echo/v4"
)

func main() {
	conf := &config.Config{}

	if err := envdecode.Decode(conf); err != nil {
		panic(fmt.Errorf("cannot read from env: %w", err))
	}

	e := echo.New()
	e.Use(middlewares.KeyAuth(conf.AuthKey))

	setupHandlers(conf, e.Group(constants.PrefixPath))

	e.Logger.Fatal(e.Start(conf.Addr()))
}

func setupHandlers(conf *config.Config, router *echo.Group) {
	// Client
	client := http.DefaultClient
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Cache.Host, conf.Cache.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Resources
	iTunesResource := resources.NewITunesResource(client)

	// Services
	songsService := services.NewSongsService([]resources.SongsResource{iTunesResource})
	songsCachedService := services.NewSongsCachedService(rdb, songsService)

	// Handlers
	songsHandler := handlers.NewSongsHandler(songsCachedService)

	// Register routes
	songsHandler.RegisterRoutes(router)
}

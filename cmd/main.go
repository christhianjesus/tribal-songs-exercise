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
	"github.com/tiaguinho/gosoap"
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
	// Clients
	client := http.DefaultClient
	clSoap, err := gosoap.SoapClient(constants.ChartLyricsPath, client)
	if err != nil {
		panic(err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", conf.Cache.Host, conf.Cache.Port),
	})

	// Resources
	iTunesResource := resources.NewITunesResource(client)
	chartLyricsResource := resources.NewChartLyricsResource(clSoap)

	// Services
	songsService := services.NewSongsService([]resources.SongsResource{
		iTunesResource,
		chartLyricsResource,
	})
	songsCachedService := services.NewSongsCachedService(rdb, songsService)

	// Handlers
	songsHandler := handlers.NewSongsHandler(songsCachedService)

	// Register routes
	songsHandler.RegisterRoutes(router)
}

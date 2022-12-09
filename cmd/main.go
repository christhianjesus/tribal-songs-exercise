package main

import (
	"christhianguevara/songs-search-exercise/cmd/config"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(conf.Addr()))
}

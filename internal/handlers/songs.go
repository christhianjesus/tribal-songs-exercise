package handlers

import (
	"christhianguevara/songs-search-exercise/domain/constants"
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type songsHandler struct {
	service services.SongsService
}

func NewSongsHandler(service services.SongsService) Handler {
	return &songsHandler{service}
}

func (h *songsHandler) RegisterRoutes(router *echo.Group) {
	router.POST(constants.SearchPath, h.Search)
}

func (h *songsHandler) Search(c echo.Context) error {
	ctx := c.Request().Context()

	params := &entities.SearchParams{}
	if err := c.Bind(params); err != nil {
		return err
	}

	songs, err := h.service.Search(ctx, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{"songs": songs}

	return c.JSON(http.StatusOK, response)
}

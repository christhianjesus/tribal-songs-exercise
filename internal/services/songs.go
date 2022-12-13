package services

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
)

type SongsService interface {
	Search(context.Context, *entities.SearchParams) ([]entities.Song, error)
}

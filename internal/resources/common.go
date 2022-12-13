package resources

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
)

type SongsResource interface {
	Search(context.Context, *entities.SearchParams) ([]entities.Song, error)
}

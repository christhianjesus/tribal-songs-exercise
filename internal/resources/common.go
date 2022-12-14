package resources

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
	"strings"
)

type SongsResource interface {
	Search(context.Context, *entities.SearchParams) ([]entities.Song, error)
}

func ContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

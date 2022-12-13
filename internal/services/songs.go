package services

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/resources"
	"context"

	"golang.org/x/sync/errgroup"
)

type SongsService interface {
	Search(context.Context, *entities.SearchParams) ([]entities.Song, error)
}

type songsService struct {
	resources []resources.SongsResource
}

func NewSongsService(resources []resources.SongsResource) SongsService {
	return &songsService{resources}
}

func (s *songsService) Search(ctx context.Context, params *entities.SearchParams) ([]entities.Song, error) {
	g, ctx := errgroup.WithContext(ctx)

	results := make([][]entities.Song, len(s.resources))
	for i, r := range s.resources {
		i, r := i, r

		g.Go(func() error {
			result, err := r.Search(ctx, params)
			if err == nil {
				results[i] = result
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return s.concat(results), nil
}

func (s *songsService) concat(results [][]entities.Song) []entities.Song {
	quantity := 0
	for _, r := range results {
		quantity += len(r)
	}

	songs := make([]entities.Song, 0, quantity)
	for _, r := range results {
		songs = append(songs, r...)
	}

	return songs
}

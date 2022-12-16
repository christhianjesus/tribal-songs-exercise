package services

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
	"encoding/json"
	"fmt"
)

type songsCachedService struct {
	cache   CacheClient
	service SongsService
}

func NewSongsCachedService(cache CacheClient, service SongsService) SongsService {
	return &songsCachedService{cache, service}
}

func (s *songsCachedService) Search(ctx context.Context, params *entities.SearchParams) ([]entities.Song, error) {
	key := s.Key(params)

	songs, err := s.Get(ctx, key)
	if err != nil {
		songs, err = s.service.Search(ctx, params)
		if err != nil {
			return nil, err
		}

		_ = s.Set(ctx, key, songs)

		return songs, nil
	}

	return songs, nil
}

func (s *songsCachedService) Key(params *entities.SearchParams) string {
	return fmt.Sprintf("%s:%s:%s", params.Artist, params.Album, params.Name)
}

func (s *songsCachedService) Get(ctx context.Context, key string) ([]entities.Song, error) {
	val, err := s.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var songs []entities.Song
	if err = json.Unmarshal(val, &songs); err != nil {
		return nil, err
	}

	return songs, nil
}

func (s *songsCachedService) Set(ctx context.Context, key string, songs []entities.Song) error {
	val, err := json.Marshal(songs)
	if err != nil {
		return err
	}

	return s.cache.Set(ctx, key, string(val), 0).Err()
}

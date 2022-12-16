package resources

import (
	"christhianguevara/songs-search-exercise/domain/constants"
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
)

type chartLyricsResource struct {
	client SoapClient
}

func NewChartLyricsResource(client SoapClient) SongsResource {
	return &chartLyricsResource{client}
}

func (c *chartLyricsResource) Search(ctx context.Context, params *entities.SearchParams) ([]entities.Song, error) {
	if params.Name == "" || params.Artist == "" {
		return nil, nil
	}

	res, err := c.client.Call("SearchLyric", c.buildQueryParams(params))
	if err != nil {
		return nil, err
	}

	var response entities.SearchLyricResponse
	if err = res.Unmarshal(&response); err != nil {
		return nil, err
	}

	songs := make([]entities.Song, 0, len(response.Results))
	for _, song := range response.Results {
		if song.ID != 0 {
			song.Origin = constants.ChartLyricsOrigin
			songs = append(songs, entities.Song(song))
		}
	}

	return songs, nil
}

func (c *chartLyricsResource) buildQueryParams(params *entities.SearchParams) map[string]interface{} {
	return map[string]interface{}{
		"song":   params.Name,
		"artist": params.Artist,
	}
}

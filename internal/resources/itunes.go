package resources

import (
	"christhianguevara/songs-search-exercise/domain/constants"
	"christhianguevara/songs-search-exercise/domain/entities"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type iTunesResource struct {
	client HTTPClient
}

func NewITunesResource(client HTTPClient) SongsResource {
	return &iTunesResource{client}
}

func (i *iTunesResource) Search(ctx context.Context, params *entities.SearchParams) ([]entities.Song, error) {
	if params.Name == "" && params.Album == "" && params.Artist == "" {
		return nil, nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", constants.ITunesPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.URL.RawQuery = i.buildQueryParams(params)

	var resp *http.Response
	if resp, err = i.client.Do(req); err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	var response entities.ITunesResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	songs := make([]entities.Song, 0, len(response.Results))
	for _, song := range response.Results {
		song.Origin = constants.ITunesOrigin
		songs = append(songs, entities.Song(song))
	}

	return i.filterResults(params, songs), nil
}

func (i *iTunesResource) buildQueryParams(params *entities.SearchParams) string {
	var term, attribute string
	switch {
	case params.Name != "":
		term = params.Name
		attribute = "songTerm"
	case params.Album != "":
		term = params.Album
		attribute = "albumTerm"
	case params.Artist != "":
		term = params.Artist
		attribute = "artistTerm"
	}

	q := url.Values{}
	q.Add("term", term)
	q.Add("attribute", attribute)
	q.Add("media", "music")
	q.Add("entity", "song")
	q.Add("limit", "200")

	return q.Encode()
}

func (i *iTunesResource) filterResults(params *entities.SearchParams, songs []entities.Song) []entities.Song {
	filters := *params
	switch {
	case filters.Name != "":
		// Do nothing
	case filters.Album != "":
		filters.Album = ""
	case filters.Artist != "":
		filters.Artist = ""
	}

	filteredSongs := make([]entities.Song, 0, len(songs))
	for _, s := range songs {
		c1 := ContainsI(s.Album, filters.Album)
		c2 := ContainsI(s.Artist, filters.Artist)
		if c1 && c2 {
			filteredSongs = append(filteredSongs, s)
		}
	}

	return filteredSongs
}

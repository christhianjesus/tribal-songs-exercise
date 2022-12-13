package services

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/resources"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type songsServiceMock struct {
	resources []resources.SongsResource
	service   SongsService
}

func Test_SearchService(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		expectError bool
		expectMSG   string
		expectLen   int
		function    func(*songsServiceMock)
	}{
		{
			name:        "Error all apis failed",
			expectError: true,
			expectMSG:   "any error",
			expectLen:   0,
			function:    func(s *songsServiceMock) { s.expectSearchErrors() },
		},
		{
			name:        "Error the apis partially responded",
			expectError: true,
			expectMSG:   "any error",
			expectLen:   0,
			function:    func(s *songsServiceMock) { s.expectSearchErrorAndOK() },
		},
		{
			name:        "Songs Search OK",
			expectError: false,
			expectLen:   3,
			function:    func(s *songsServiceMock) { s.expectSearchOK() },
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			r := setupSongsService()
			sc.function(r)

			songs, err := r.service.Search(context.Background(), nil)

			if sc.expectError {
				assert.Error(t, err)
				assert.Equal(t, sc.expectMSG, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Len(t, songs, sc.expectLen)
		})
	}

}

func setupSongsService() *songsServiceMock {
	resourcesMock := []resources.SongsResource{
		new(mocks.SongsResource),
		new(mocks.SongsResource),
	}

	return &songsServiceMock{
		resources: resourcesMock,
		service:   NewSongsService(resourcesMock),
	}
}

func (s *songsServiceMock) expectSearchErrors() {
	s.resources[0].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))

	s.resources[1].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))
}

func (s *songsServiceMock) expectSearchErrorAndOK() {
	s.resources[0].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))

	s.resources[1].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{{}}, nil)
}

func (s *songsServiceMock) expectSearchOK() {
	s.resources[0].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{{}}, nil)

	s.resources[1].(*mocks.SongsResource).
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{{}, {}}, nil)
}

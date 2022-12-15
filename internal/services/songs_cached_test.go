package services

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_SearchCachedService(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		expectError bool
		expectMSG   string
		expectLen   int
		function    func(*songsCachedServiceMock)
	}{
		{
			name:        "Get cache OK",
			expectError: false,
			expectLen:   1,
			function:    func(s *songsCachedServiceMock) { s.expectGetOK() },
		},
		{
			name:        "Get cache Error and API Error",
			expectError: true,
			expectMSG:   "any error",
			function:    func(s *songsCachedServiceMock) { s.expectGetError(); s.expectSearchError() },
		},
		{
			name:        "Get cache Error and API OK and Set cache Error",
			expectError: false,
			expectLen:   1,
			function:    func(s *songsCachedServiceMock) { s.expectGetError(); s.expectSearchOK(); s.expectSetError() },
		},
		{
			name:        "Get cache Error and API OK and Set cache Ok",
			expectError: false,
			expectLen:   1,
			function:    func(s *songsCachedServiceMock) { s.expectGetError(); s.expectSearchOK(); s.expectSetOK() },
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			r := setupSongsCachedService(t)
			sc.function(r)

			songs, err := r.cachedService.Search(context.Background(), &entities.SearchParams{})

			if sc.expectError {
				assert.Error(t, err)
				assert.Equal(t, sc.expectMSG, err.Error())
				assert.Nil(t, songs)
			} else {
				assert.NoError(t, err)
				assert.Len(t, songs, sc.expectLen)
			}
		})
	}
}

type songsCachedServiceMock struct {
	redis         *mocks.CacheClient
	service       *mocks.SongsService
	cachedService SongsService
}

func setupSongsCachedService(t *testing.T) *songsCachedServiceMock {
	client := mocks.NewCacheClient(t)
	serviceMock := mocks.NewSongsService(t)

	return &songsCachedServiceMock{
		redis:         client,
		service:       serviceMock,
		cachedService: NewSongsCachedService(client, serviceMock),
	}
}

func (h *songsCachedServiceMock) expectSearchError() {
	h.service.
		On("Search", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))
}

func (h *songsCachedServiceMock) expectSearchOK() {
	h.service.
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{{}}, nil)
}

func (h *songsCachedServiceMock) expectGetError() {
	result := &redis.StringCmd{}
	result.SetErr(errors.New("any error"))

	h.redis.
		On("Get", mock.Anything, mock.Anything).
		Return(result)
}

func (h *songsCachedServiceMock) expectGetOK() {
	result := &redis.StringCmd{}
	result.SetVal(`[{}]`)

	h.redis.
		On("Get", mock.Anything, mock.Anything).
		Return(result)
}

func (h *songsCachedServiceMock) expectSetError() {
	result := &redis.StatusCmd{}
	result.SetErr(errors.New("any error"))

	h.redis.
		On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(result)
}

func (h *songsCachedServiceMock) expectSetOK() {
	result := &redis.StatusCmd{}

	h.redis.
		On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(result)
}

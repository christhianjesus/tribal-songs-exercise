package resources

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type iTunesResourceMock struct {
	httpc    *mocks.HTTPClient
	resource SongsResource
}

func Test_iTunesResource(t *testing.T) {
	t.Parallel()

	expectedResponse := []entities.Song{
		{
			ID:       1469577741,
			Artist:   "Jack Johnson",
			Album:    "Jack Johnson and Friends: Sing-A-Longs and Lullabies for the Film Curious George",
			Name:     "Upside Down",
			Price:    1.29,
			Duration: 208643,
			Origin:   "apple",
		},
	}

	cases := []struct {
		name           string
		params         *entities.SearchParams
		expectError    bool
		expectMSG      string
		expectResponse []entities.Song
		ctx            context.Context
		function       func(*iTunesResourceMock)
	}{
		{
			name:           "Empty params",
			params:         &entities.SearchParams{},
			expectError:    false,
			expectResponse: nil,
			ctx:            context.Background(),
			function:       func(s *iTunesResourceMock) {},
		},
		{
			name:        "Error NewRequestWithContext fail",
			params:      &entities.SearchParams{Name: "some name"},
			expectError: true,
			expectMSG:   "net/http: nil Context",
			ctx:         nil,
			function:    func(s *iTunesResourceMock) {},
		},
		{
			name:        "Error all apis failed",
			params:      &entities.SearchParams{Name: "some name"},
			expectError: true,
			expectMSG:   "any error",
			ctx:         context.Background(),
			function:    func(s *iTunesResourceMock) { s.expectAPIError() },
		},
		{
			name:        "Error read error",
			params:      &entities.SearchParams{Name: "some name"},
			expectError: true,
			expectMSG:   "any error",
			ctx:         context.Background(),
			function:    func(s *iTunesResourceMock) { s.expectAPIConnectionFails() },
		},
		{
			name:        "Error api empty response",
			params:      &entities.SearchParams{Name: "some name"},
			expectError: true,
			expectMSG:   "unexpected end of JSON input",
			ctx:         context.Background(),
			function:    func(s *iTunesResourceMock) { s.expectAPIEmptyResponse() },
		},
		{
			name:           "API response OK",
			params:         &entities.SearchParams{Name: "some name"},
			expectError:    false,
			expectResponse: expectedResponse,
			ctx:            context.Background(),
			function:       func(s *iTunesResourceMock) { s.expectAPIResponseOK() },
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			r := setupITunesResource(t)
			sc.function(r)

			songs, err := r.resource.Search(sc.ctx, sc.params)

			if sc.expectError {
				assert.Error(t, err)
				assert.Equal(t, sc.expectMSG, err.Error())
				assert.Nil(t, songs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, sc.expectResponse, songs)
			}
		})
	}

}

func setupITunesResource(t *testing.T) *iTunesResourceMock {
	client := mocks.NewHTTPClient(t)

	return &iTunesResourceMock{
		httpc:    client,
		resource: NewITunesResource(client),
	}
}

func (r *iTunesResourceMock) expectAPIError() {
	r.httpc.
		On("Do", mock.Anything).
		Return(nil, errors.New("any error"))
}

func (r *iTunesResourceMock) expectAPIConnectionFails() {
	r.httpc.
		On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(&FailRead{}),
		}, nil)
}

func (r *iTunesResourceMock) expectAPIEmptyResponse() {
	r.httpc.
		On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(``)),
		}, nil)
}

func (r *iTunesResourceMock) expectAPIResponseOK() {
	apiResponse := `{"results": [{
		"trackId":1469577741,
		"artistName":"Jack Johnson",
		"collectionName":"Jack Johnson and Friends: Sing-A-Longs and Lullabies for the Film Curious George",
		"trackName":"Upside Down",
		"trackPrice":1.29,
		"trackTimeMillis":208643
	}]}`

	r.httpc.
		On("Do", mock.Anything).
		Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(apiResponse)),
		}, nil)
}

type FailRead struct{}

func (*FailRead) Read(p []byte) (n int, err error) {
	return 0, errors.New("any error")
}

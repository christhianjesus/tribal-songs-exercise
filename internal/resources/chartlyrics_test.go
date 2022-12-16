package resources

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tiaguinho/gosoap"
)

type chartLyricsResourceMock struct {
	sc       *mocks.SoapClient
	resource SongsResource
}

func Test_ChartLyricsResource(t *testing.T) {
	t.Parallel()

	expectedResponse := []entities.Song{
		{
			ID:     6133139,
			Name:   "Ahora Quien",
			Artist: "Marc Anthony",
			Origin: "chartlyrics",
		},
	}

	cases := []struct {
		name           string
		params         *entities.SearchParams
		expectError    bool
		expectMSG      string
		expectResponse []entities.Song
		function       func(*chartLyricsResourceMock)
	}{
		{
			name:           "Empty params",
			params:         &entities.SearchParams{},
			expectError:    false,
			expectResponse: nil,
			function:       func(s *chartLyricsResourceMock) {},
		},
		{
			name:        "Error all apis failed",
			params:      &entities.SearchParams{Name: "some name", Artist: "some artist"},
			expectError: true,
			expectMSG:   "any error",
			function:    func(s *chartLyricsResourceMock) { s.expectAPIError() },
		},
		{
			name:        "Error api empty response",
			params:      &entities.SearchParams{Name: "some name", Artist: "some artist"},
			expectError: true,
			expectMSG:   "Body is empty",
			function:    func(s *chartLyricsResourceMock) { s.expectAPIEmptyResponse() },
		},
		{
			name:           "API response OK",
			params:         &entities.SearchParams{Name: "some name", Artist: "some artist"},
			expectError:    false,
			expectResponse: expectedResponse,
			function:       func(s *chartLyricsResourceMock) { s.expectAPIResponseOK() },
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			r := setupChartLyricsResource(t)
			sc.function(r)

			songs, err := r.resource.Search(context.Background(), sc.params)

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

func setupChartLyricsResource(t *testing.T) *chartLyricsResourceMock {
	soapClient := mocks.NewSoapClient(t)

	return &chartLyricsResourceMock{
		sc:       soapClient,
		resource: NewChartLyricsResource(soapClient),
	}
}

func (r *chartLyricsResourceMock) expectAPIError() {
	r.sc.
		On("Call", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))
}

func (r *chartLyricsResourceMock) expectAPIEmptyResponse() {
	result := &gosoap.Response{}

	r.sc.
		On("Call", mock.Anything, mock.Anything).
		Return(result, nil)
}

func (r *chartLyricsResourceMock) expectAPIResponseOK() {
	apiResponse := `<SearchLyricResponse xmlns="http://api.chartlyrics.com/">
		<SearchLyricResult>
			<SearchLyricResult>
				<TrackId>6133139</TrackId>
				<Artist>Marc Anthony</Artist>
				<Song>Ahora Quien</Song>
			</SearchLyricResult>
		</SearchLyricResult>
		<SearchLyricResult xsi:nil="true" />
	</SearchLyricResponse>`

	result := &gosoap.Response{
		Body: []byte(apiResponse),
	}

	r.sc.
		On("Call", mock.Anything, mock.Anything).
		Return(result, nil)
}

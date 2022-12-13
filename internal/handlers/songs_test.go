package handlers

import (
	"christhianguevara/songs-search-exercise/domain/entities"
	"christhianguevara/songs-search-exercise/internal/tests/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type songsHandler struct {
	service *mocks.SongsService
	handler Handler
}

func Test_SearchHandler(t *testing.T) {
	t.Parallel()

	var (
		paramsError = `{"name": 1}`
		paramsOk    = `{"name": "some"}`
	)

	cases := []struct {
		name         string
		params       string
		expectedCode int
		expectedBody string
		function     func(*songsHandler)
	}{
		{
			name:         "Error invalid params",
			params:       paramsError,
			expectedCode: http.StatusBadRequest,
			expectedBody: "{\"message\":\"Unmarshal type error: expected=string, got=number, field=name, offset=10\"}\n",
			function:     func(sh *songsHandler) {},
		},
		{
			name:         "Error api response",
			params:       paramsOk,
			expectedCode: http.StatusInternalServerError,
			expectedBody: "{\"message\":\"any error\"}\n",
			function:     func(sh *songsHandler) { sh.expectSearchError() },
		},
		{
			name:         "Songs Search OK",
			params:       paramsOk,
			expectedCode: http.StatusOK,
			expectedBody: "{\"songs\":[]}\n",
			function:     func(sh *songsHandler) { sh.expectSearchOK() },
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			r := setupSongsHandler()
			sc.function(r)

			e := echo.New()
			e.POST("/search", r.handler.(*SongsHandler).Search)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/search", strings.NewReader(sc.params))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(rec, req)

			assert.Equal(t, sc.expectedCode, rec.Code)
			assert.Equal(t, sc.expectedBody, rec.Body.String())
		})
	}

}

func setupSongsHandler() *songsHandler {
	serviceMock := new(mocks.SongsService)

	return &songsHandler{
		service: serviceMock,
		handler: NewSongsHandler(serviceMock),
	}
}

func (h *songsHandler) expectSearchError() {
	h.service.
		On("Search", mock.Anything, mock.Anything).
		Return(nil, errors.New("any error"))
}

func (h *songsHandler) expectSearchOK() {
	h.service.
		On("Search", mock.Anything, mock.Anything).
		Return([]entities.Song{}, nil)
}

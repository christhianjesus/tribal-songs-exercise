package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_AdminAuth(t *testing.T) {
	t.Parallel()

	const exampleToken = "123"

	cases := []struct {
		name         string
		header       string
		expectedCode int
	}{
		{
			name:         "No token",
			header:       "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid token",
			header:       "invalid",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "OK",
			header:       exampleToken,
			expectedCode: http.StatusOK,
		},
	}

	for _, sc := range cases {
		t.Run(sc.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/", func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			})
			e.Use(KeyAuth(exampleToken))

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", sc.header)
			res := httptest.NewRecorder()

			e.ServeHTTP(res, req)

			assert.Equal(t, sc.expectedCode, res.Code)
		})
	}
}

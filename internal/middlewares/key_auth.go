package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func KeyAuth(authKey string) echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestKey := c.Request().Header.Get("Authorization")

			switch requestKey {
			case authKey:
				return handler(c)
			case "":
				return c.NoContent(http.StatusBadRequest)
			default:
				return c.NoContent(http.StatusUnauthorized)
			}
		}
	}
}

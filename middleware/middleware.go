package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u, p, ok := c.Request().BasicAuth()
		if !ok || (u != os.Getenv("ADMIN_USERNAME") || p != os.Getenv("ADMIN_PASSWORD")) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"message": http.StatusText(http.StatusUnauthorized),
			})
		}

		return next(c)
	}
}

package middleware

import (
	"net/http"

	"github.com/Gergenus/StandardLib/pkg"
	"github.com/labstack/echo/v4"
)

type Middleware interface {
	AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type EchoMiddleware struct {
	JWT pkg.JWTpkg
}

func NewEchoMiddleware(JWT pkg.JWTpkg) EchoMiddleware {
	return EchoMiddleware{JWT: JWT}
}

func (e *EchoMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("Auth")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "No auth token",
			})
		}
		token := cookie.Value
		name, err := e.JWT.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}
		c.Set("name", name)
		return next(c)
	}
}

func (e *EchoMiddleware) WSAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		name, err := e.JWT.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid token",
			})
		}
		c.Set("name", name)
		return next(c)
	}
}

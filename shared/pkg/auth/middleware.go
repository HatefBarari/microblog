package auth

import (
	"net/http"
	"strings"

	"github.com/HatefBarari/microblog-shared/pkg/httputil"
	"github.com/labstack/echo/v4"
)

func Middleware(accessSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httputil.NewError(401, "missing token"))
			}
			bearer := strings.Split(authHeader, " ")
			if len(bearer) != 2 || bearer[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httputil.NewError(401, "invalid token format"))
			}
			claims, err := ValidateToken(bearer[1], accessSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httputil.NewError(401, "invalid or expired token"))
			}
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}
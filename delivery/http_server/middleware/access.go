package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mohamadrezamomeni/momo/service/auth"
)

func AccessCheck(auth *auth.Auth, isAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.ErrUnauthorized
			}

			clam, isValid, err := auth.DecodeToken(authHeader)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error":   true,
					"message": "Internal Server Error: something went wrong",
				})
			}
			if !isValid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   true,
					"message": "user isn't authorized",
				})
			}
			if isAdmin && clam.IsAdmin != true {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error":   true,
					"message": "user isn't allowed to call this method",
				})
			}
			return next(c)
		}
	}
}

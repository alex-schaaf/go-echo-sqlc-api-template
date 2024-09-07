package auth

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CookieAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		if isValid := IsValidToken(cookie.Value); !isValid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}
		payload, err := GetTokenPayload(cookie.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		sub, ok := payload["sub"].(float64)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		userID := int64(sub)
		if c.Param("user_id") == "" { // any authorized user can access non-user scoped routes
			c.Set("user_id", userID)
			return next(c)
		}
		// user scoped routes, check if the requesting user is the same as the user in the route params
		if c.Param("user_id") != strconv.FormatInt(userID, 10) {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden")
		}
		c.Set("user_id", userID)
		return next(c)
	}
}

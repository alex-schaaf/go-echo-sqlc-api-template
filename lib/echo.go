package lib

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// GetEchoInstance with middlewares
func GetEchoInstance() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

// GetUserID returns the user ID from the echo context.
// This function is used in the API handlers to get the user ID from the context.
// The user ID is set in the CookieAuthMiddleware.
func GetUserID(e echo.Context) int64 {
	userIDStr := e.Get("user_id").(string)
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)
	return userID
}

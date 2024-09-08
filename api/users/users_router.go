package users

import (
	"app/db"
	"app/lib/auth"
	"app/lib/config"

	"github.com/labstack/echo/v4"
)

func AddUsersRouter(e *echo.Echo, queries *db.Queries, config *config.Config) {
	usersHandler := NewUsersHandler(queries, config)
	usersGroup := e.Group("/users")
	usersGroup.Use(auth.CookieAuthMiddleware(config.JWT_SECRET))
	usersGroup.PATCH("/:user_id/password", usersHandler.UpdateUserPasswordHandler)
	usersGroup.DELETE("/:user_id", usersHandler.DeleteUserHandler)
}

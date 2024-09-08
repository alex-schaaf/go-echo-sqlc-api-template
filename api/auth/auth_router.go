package auth

import (
	"app/db"
	"app/lib/config"

	"github.com/labstack/echo/v4"
)

func AddAuthRouter(e *echo.Echo, queries *db.Queries, config *config.Config) {
	authHandler := NewAuthHandler(queries, config)
	usersGroup := e.Group("/auth")
	usersGroup.POST("/login", authHandler.LoginHandler)
	usersGroup.POST("/logout", authHandler.LogoutHandler)
	usersGroup.POST("/register", authHandler.RegisterHandler)
}

package main

import (
	"app/api/auth"
	"app/api/users"
	"app/db"
	"app/lib"
	libAuth "app/lib/auth"
	"app/lib/config"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func initializeDatabase(config *config.Config) *sql.DB {
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		panic(err)
	}
	return db
}

func getQueries(conn *sql.DB) *db.Queries {
	return db.New(conn)
}

func main() {
	config := config.InitConfig()
	conn := initializeDatabase(config)
	queries := getQueries(conn)

	e := lib.GetEchoInstance()

	authHandler := auth.NewAuthHandler(queries)

	e.POST("/auth/login", authHandler.LoginHandler)
	e.POST("/auth/logout", authHandler.LogoutHandler)
	e.POST("/auth/register", authHandler.RegisterHandler)

	usersHandler := users.NewUsersHandler(queries)
	usersGroup := e.Group("/users")
	usersGroup.Use(libAuth.CookieAuthMiddleware)
	usersGroup.PATCH("/:user_id/password", usersHandler.UpdateUserPasswordHandler)
	usersGroup.DELETE("/:user_id", usersHandler.DeleteUserHandler)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))
}

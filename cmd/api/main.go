package main

import (
	"app/api/auth"
	"app/api/users"
	"app/db"
	"app/lib"
	"app/lib/config"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func getQueries(conn *sql.DB) *db.Queries {
	return db.New(conn)
}

func main() {
	config := config.InitConfig()
	conn := lib.InitializeDatabase(config)
	queries := getQueries(conn)

	e := lib.GetEchoInstance()

	auth.AddAuthRouter(e, queries, config)
	users.AddUsersRouter(e, queries, config)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))
}

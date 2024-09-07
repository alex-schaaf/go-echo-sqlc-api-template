package test

import (
	"app/db"
	"app/lib"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
)

func MigrateDatabase(conn *sql.DB, migrationsPath string) error {
	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	if err != nil {
		return err
	}
	sourceUrl := fmt.Sprintf("file://%s", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(sourceUrl, "sqlite3", driver)
	if err != nil {
		return err
	}
	return m.Up()
}

func InitializeMemoryDatabase() *sql.DB {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	if err := MigrateDatabase(conn, "../../db/migrations"); err != nil {
		panic(err)
	}
	return conn
}

func SetupTest() (*echo.Echo, *db.Queries) {
	conn := InitializeMemoryDatabase()
	queries := db.New(conn)
	return lib.GetEchoInstance(), queries
}

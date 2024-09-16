package lib

import (
	"app/lib/config"
	"database/sql"
)

func InitializeDatabase(config *config.Config) *sql.DB {
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		panic(err)
	}
	return db
}

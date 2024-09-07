package main

import (
	"app/lib/config"
	"app/lib/test"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config := config.InitConfig()
	// Initialize the database connection
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Run the migrations
	if err := test.MigrateDatabase(db, "db/migrations"); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations ran successfully")
}

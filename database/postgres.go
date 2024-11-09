package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // Importing the PostgreSQL driver for Go
)

func NewPostgresDB() *sql.DB {
	connStr := "user=postgres dbname=inventaris sslmode=disable password=postgres host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err.Error())
	}

	return db
}

package db

import (
	"database/sql"
	"fmt"
    "log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes the global DB connection
func Connect() {
	// Load connection info from .env or defaults
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping to verify connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	log.Println("Connected to the database successfully.")
}

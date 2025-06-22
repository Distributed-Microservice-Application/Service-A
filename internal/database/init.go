package db

import (
	"database/sql"
	"log"
	"os"

	// PostgreSQL driver
	_ "github.com/lib/pq"
)

const (
	database = "DB_URL"
)

func INIT_DB() *sql.DB {
	DB_URL := getEnv("DB_URL", database)

	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

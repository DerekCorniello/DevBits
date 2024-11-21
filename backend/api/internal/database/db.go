package database

import (
	"database/sql"
	"log"
)

var DB *sql.DB // Global database instance

// Connect initializes a database connection
func Connect(dsn string, driverName string) {
	var err error
	DB, err = sql.Open(driverName, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
}

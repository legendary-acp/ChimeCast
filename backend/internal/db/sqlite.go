package db

import (
	"database/sql"
	"log"
)

func CreateDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./chimecast.db")
	if err != nil {
		return nil, err
	}
	if err = createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	err := createUserTable(db)
	if err != nil {
		return err
	}
	return nil
}

func createUserTable(db *sql.DB) error {
	// Create user table
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS user (
		"UserId" TEXT PRIMARY KEY,         -- Unique UserId for the user (acts as the primary key)
		"Email" TEXT UNIQUE,               -- Unique email address for the user
		"Name" TEXT,                       -- Name of the user
		"Password" TEXT,                   -- Hashed password for authentication
		"CreatedAt" DATETIME               -- Time of creating user
	);`
	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		log.Printf("Error creating User table: %s", err)
		return err
	}
	return nil
}

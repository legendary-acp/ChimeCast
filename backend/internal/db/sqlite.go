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

	err = createRoomTable(db)
	if err != nil {
		return err
	}

	return nil
}

func createUserTable(db *sql.DB) error {
	// Create user table
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"userName" TEXT PRIMARY KEY,       -- Unique UserId for the user (acts as the primary key)
		"id" TEXT UNIQUE,                  -- Unique ID for the user
		"email" TEXT UNIQUE,               -- Unique email address for the user
		"name" TEXT,                       -- Name of the user
		"hashedPassword" TEXT,             -- Hashed password for authentication
		"createdAt" DATETIME               -- Time of creating user
	);`
	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		log.Printf("Error creating User table: %s", err)
		return err
	}
	return nil
}

func createRoomTable(db *sql.DB) error {
	// Create rooms table
	createRoomTableSQL := `CREATE TABLE IF NOT EXISTS rooms (
		"id" TEXT PRIMARY KEY,          -- Unique ID for the room
		"name" TEXT,                    -- Name of the room
		"createdAt" DATETIME,           -- Time of creating room
		"status" INTEGER                -- Room status: 0 for inactive, 1 for active
	);`
	_, err := db.Exec(createRoomTableSQL)
	if err != nil {
		log.Printf("Error creating Rooms table: %s", err)
		return err
	}
	return nil
}

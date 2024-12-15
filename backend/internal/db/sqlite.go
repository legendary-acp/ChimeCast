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
	// Create user table first since rooms will reference it
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
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
        "Username" TEXT PRIMARY KEY,    -- Unique UserId for the user (acts as the primary key)
        "ID" TEXT UNIQUE,              -- Unique ID for the user
        "Email" TEXT UNIQUE,           -- Unique email address for the user
        "Name" TEXT,                   -- Name of the user
        "HashedPassword" TEXT,         -- Hashed password for authentication
        "CreatedAt" DATETIME           -- Time of creating user
    );`

	_, err := db.Exec(createUserTableSQL)
	if err != nil {
		log.Printf("Error creating User table: %s", err)
		return err
	}
	return nil
}

func createRoomTable(db *sql.DB) error {
	createRoomTableSQL := `CREATE TABLE IF NOT EXISTS rooms (
        "ID" TEXT PRIMARY KEY,         -- Unique ID for the room
        "Name" TEXT,                   -- Name of the room
        "HostID" TEXT NOT NULL,        -- ID of the user who created the room
        "CreatedAt" DATETIME,          -- Time of creating room
        "Status" INTEGER,              -- Room status: 0 for inactive, 1 for active
        FOREIGN KEY ("HostID") REFERENCES users("ID")
    );`

	_, err := db.Exec(createRoomTableSQL)
	if err != nil {
		log.Printf("Error creating Rooms table: %s", err)
		return err
	}
	return nil
}

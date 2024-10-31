package repositories

import "database/sql"

type AuthRepository struct {
	DB *sql.DB
}

type RoomRepository struct {
	DB *sql.DB
}

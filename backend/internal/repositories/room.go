package repositories

import (
	"database/sql"
	"errors"
	"log"
	"sort"

	"github.com/legendary-acp/chimecast/internal/models"
)

func NewRoomRepositor(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		DB: db,
	}
}

func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room

	rows, err := r.DB.Query("SELECT * FROM rooms")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room

		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedAt, &room.Status); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Status != 0 && rooms[j].Status == 0
	})

	return rooms, nil
}

func (r *RoomRepository) CreateRoom(room *models.Room) error {
	_, err := r.DB.Exec("INSERT INTO rooms (ID, Name, CreatedAt, Status) VALUES (?, ?, ?, ?)", room.ID, room.Name, room.CreatedAt, room.Status)
	if err != nil {
		log.Println("unable to insert room into db: " + err.Error())
		return errors.New("unable to insert room into db: " + err.Error())
	}

	log.Println("Room created: " + room.Name)
	return nil
}

func (r *RoomRepository) DoesRoomExist(ID string) (*bool, error) {
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = ?)", ID).Scan(&exists)
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}

	return &exists, nil
}

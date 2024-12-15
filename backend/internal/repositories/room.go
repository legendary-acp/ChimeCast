package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/legendary-acp/chimecast/internal/models"
)

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{
		DB: db,
	}
}

func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	var rooms []models.Room
	rows, err := r.DB.Query(`
        SELECT ID, Name, HostID, CreatedAt, Status 
        FROM rooms
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		if err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.HostID,
			&room.CreatedAt,
			&room.Status,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Sort active rooms first
	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Status == models.RoomStatusActive && rooms[j].Status != models.RoomStatusActive
	})

	return rooms, nil
}

func (r *RoomRepository) CreateRoom(room *models.Room) error {
	_, err := r.DB.Exec(`
        INSERT INTO rooms (ID, Name, HostID, CreatedAt, Status) 
        VALUES (?, ?, ?, ?, ?)`,
		room.ID,
		room.Name,
		room.HostID,
		room.CreatedAt,
		room.Status,
	)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return fmt.Errorf("failed to create room: %v", err)
	}

	log.Printf("Room created: %s (ID: %s)", room.Name, room.ID)
	return nil
}

func (r *RoomRepository) GetRoom(roomID string) (*models.Room, error) {
	var room models.Room
	err := r.DB.QueryRow(`
       	SELECT ID, Name, HostID, CreatedAt, Status 
        FROM rooms 
        WHERE id = ?`,
		roomID,
	).Scan(
		&room.ID,
		&room.Name,
		&room.HostID,
		&room.CreatedAt,
		&room.Status,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("room not found")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &room, nil
}

func (r *RoomRepository) DoesRoomExist(ID string) (*bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 
            FROM rooms 
            WHERE id = ?
        )`,
		ID,
	).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &exists, nil
}

func (r *RoomRepository) UpdateRoomStatus(roomID string, status int) error {
	result, err := r.DB.Exec(`
        UPDATE rooms 
        SET status = ? 
        WHERE id = ?`,
		status,
		roomID,
	)
	if err != nil {
		return fmt.Errorf("failed to update room status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking update result: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("room not found")
	}

	return nil
}

func (r *RoomRepository) DeleteRoom(roomID string) error {
	result, err := r.DB.Exec(`
        DELETE FROM rooms 
        WHERE id = ?`,
		roomID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete room: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking delete result: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("room not found")
	}

	return nil
}

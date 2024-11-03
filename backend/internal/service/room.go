package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/utils"
)

func NewRoomService(roomRepository *repositories.RoomRepository) *RoomService {
	return &RoomService{
		RoomRepository: roomRepository,
		Connections:    make(map[string][]*websocket.Conn),
	}
}

func (r *RoomService) GetAllRooms() ([]models.Room, error) {
	return r.RoomRepository.GetAllRooms()
}

func (r *RoomService) CreateRoom(request *models.CreateRoomRequest) (*string, error) {
	if request.Name == "" {
		return nil, errors.New("name can't be empty")
	}

	var room models.Room

	room.Name = request.Name
	room.ID = utils.CreateNewUUID()
	room.Status = 1
	room.CreatedAt = time.Now()

	if err := r.RoomRepository.CreateRoom(&room); err != nil {
		return nil, err
	}
	return &room.ID, nil
}

func (r *RoomService) JoinRoom(ID string) (*bool, error) {
	return r.RoomRepository.DoesRoomExist(ID)
}

func (r *RoomService) WebSocketConnection(roomID string, conn *websocket.Conn) error {
	exists, err := r.RoomRepository.DoesRoomExist(roomID)
	if err != nil || !*exists {
		log.Printf("Room with ID %s doesn't exist: %v", roomID, err)
		return fmt.Errorf("room does not exist: %s", roomID)
	}

	r.Connections[roomID] = append(r.Connections[roomID], conn)
	log.Printf("New connection added to room %s. Total connections: %d", roomID, len(r.Connections[roomID]))

	defer func() {
		r.removeConnection(roomID, conn)
		log.Printf("WebSocket connection closed for room %s", roomID)
		conn.Close()
	}()

	for {
		var message map[string]interface{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Printf("Error reading message from room %s: %v", roomID, err)
			return fmt.Errorf("error reading message: %v", err)
		}

		if err := r.broadcastToRoom(roomID, message, conn); err != nil {
			log.Printf("Error broadcasting message in room %s: %v", roomID, err)
			return err
		}
	}
}

func (r *RoomService) broadcastToRoom(roomID string, message map[string]interface{}, sender *websocket.Conn) error {
	for _, conn := range r.Connections[roomID] {
		if conn != sender {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to a client in room %s: %v", roomID, err)
				return err
			}
		}
	}
	return nil
}

func (r *RoomService) removeConnection(roomID string, conn *websocket.Conn) {
	connections := r.Connections[roomID]
	for i, c := range connections {
		if c == conn {
			r.Connections[roomID] = append(connections[:i], connections[i+1:]...)
			log.Printf("Connection removed from room %s. Total connections: %d", roomID, len(r.Connections[roomID]))
			break
		}
	}
}

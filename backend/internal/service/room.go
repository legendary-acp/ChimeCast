package service

import (
	"errors"
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

func (r *RoomService) WebSocketConnection(roomID string, conn *websocket.Conn) {
	exists, err := r.RoomRepository.DoesRoomExist(roomID)
	if err != nil || !*exists {
		log.Printf("room with id %s doesn't exist", roomID)
		return
	}

	r.Connections[roomID] = append(r.Connections[roomID], conn)

	for {
		var message map[string]interface{}

		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("Error reading message: " + err.Error())
			break
		}

		r.broadcastToRoom(roomID, message, conn)
	}

	r.removeConnection(roomID, conn)
}

// Broadcasts a message to all connections in the room except the sender
func (r *RoomService) broadcastToRoom(roomID string, message map[string]interface{}, sender *websocket.Conn) {
	for _, conn := range r.Connections[roomID] {
		if conn != sender {
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}
		}
	}
}

// Removes a connection from the room
func (r *RoomService) removeConnection(roomID string, conn *websocket.Conn) {
	connections := r.Connections[roomID]
	for i, c := range connections {
		if c == conn {
			// Remove the connection by slicing the slice
			r.Connections[roomID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

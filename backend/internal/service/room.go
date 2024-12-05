package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/utils"
)

// WebRTCMessage represents a WebRTC signaling message
type WebRTCMessage struct {
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
	SenderID string          `json:"senderId"`
}

// Connection represents a WebSocket connection with user information
type Connection struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
}

// NewRoomService creates a new instance of RoomService
func NewRoomService(roomRepository *repositories.RoomRepository) *RoomService {
	return &RoomService{
		RoomRepository: roomRepository,
		Connections:    make(map[string][]*Connection),
	}
}

// GetAllRooms returns all available rooms
func (r *RoomService) GetAllRooms() ([]models.Room, error) {
	return r.RoomRepository.GetAllRooms()
}

// CreateRoom creates a new room with the given request
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

// JoinRoom checks if a room exists and can be joined
func (r *RoomService) JoinRoom(ID string) (*bool, error) {
	return r.RoomRepository.DoesRoomExist(ID)
}

// WebSocketConnection handles WebSocket connections for a room
func (r *RoomService) WebSocketConnection(roomID string, userID string, conn *websocket.Conn) error {
	exists, err := r.RoomRepository.DoesRoomExist(roomID)
	if err != nil || !*exists {
		log.Printf("Room with ID %s doesn't exist: %v", roomID, err)
		return fmt.Errorf("room does not exist: %s", roomID)
	}

	connection := &Connection{
		Conn:   conn,
		UserID: userID,
	}

	r.mu.Lock()
	if r.Connections[roomID] == nil {
		r.Connections[roomID] = make([]*Connection, 0)
	}
	r.Connections[roomID] = append(r.Connections[roomID], connection)
	r.mu.Unlock()

	log.Printf("New connection added to room %s. Total connections: %d", roomID, len(r.Connections[roomID]))

	// Notify others about new peer
	r.broadcastToRoom(roomID, WebRTCMessage{
		Type:     "user-joined",
		SenderID: userID,
	}, conn)

	defer func() {
		r.removeConnection(roomID, conn)
		log.Printf("WebSocket connection closed for room %s", roomID)
		conn.Close()
	}()

	for {
		var msg WebRTCMessage
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message from room %s: %v", roomID, err)
			return fmt.Errorf("error reading message: %v", err)
		}

		msg.SenderID = userID // Ensure sender ID is set correctly

		switch msg.Type {
		case "offer", "answer", "ice-candidate":
			// Forward these messages directly to peers
			if err := r.broadcastToRoom(roomID, msg, conn); err != nil {
				log.Printf("Error broadcasting message in room %s: %v", roomID, err)
				return err
			}

		case "leave":
			r.removeConnection(roomID, conn)
			r.broadcastToRoom(roomID, WebRTCMessage{
				Type:     "user-left",
				SenderID: userID,
			}, conn)
			return nil

		default:
			log.Printf("Unknown message type received: %s", msg.Type)
		}
	}
}

// broadcastToRoom sends a message to all connections in a room except the sender
func (r *RoomService) broadcastToRoom(roomID string, message interface{}, sender *websocket.Conn) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, conn := range r.Connections[roomID] {
		if conn.Conn != sender {
			err := conn.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to a client in room %s: %v", roomID, err)
				return err
			}
		}
	}
	return nil
}

// removeConnection removes a WebSocket connection from a room
func (r *RoomService) removeConnection(roomID string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	connections := r.Connections[roomID]
	for i, c := range connections {
		if c.Conn == conn {
			r.Connections[roomID] = append(connections[:i], connections[i+1:]...)
			log.Printf("Connection removed from room %s. Total connections: %d", roomID, len(r.Connections[roomID]))
			break
		}
	}

	// Clean up room if empty
	if len(r.Connections[roomID]) == 0 {
		delete(r.Connections, roomID)
	}
}

// GetRoomParticipants returns a list of user IDs in a room
func (r *RoomService) GetRoomParticipants(roomID string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var participants []string
	for _, conn := range r.Connections[roomID] {
		participants = append(participants, conn.UserID)
	}
	return participants
}

// IsUserInRoom checks if a user is currently in a room
func (r *RoomService) IsUserInRoom(roomID, userID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, conn := range r.Connections[roomID] {
		if conn.UserID == userID {
			return true
		}
	}
	return false
}

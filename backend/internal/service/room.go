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

type Connection struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
	JoinedAt time.Time
	Status   string // "waiting" or "admitted"
}

func NewRoomService(roomRepository *repositories.RoomRepository) *RoomService {
	return &RoomService{
		RoomRepository: roomRepository,
		Connections:    make(map[string]map[string]*Connection),
		WaitingRoom:    make(map[string]map[string]*Connection),
	}
}

// GetAllRooms returns all available rooms
func (r *RoomService) GetAllRooms() ([]models.Room, error) {
	return r.RoomRepository.GetAllRooms()
}

func (r *RoomService) CreateRoom(request *models.CreateRoomRequest, hostID string) (*string, error) {
	if request.Name == "" {
		return nil, errors.New("name can't be empty")
	}

	room := models.Room{
		ID:        utils.CreateNewUUID(),
		Name:      request.Name,
		HostID:    hostID,
		Status:    models.RoomStatusActive,
		CreatedAt: time.Now(),
	}

	if err := r.RoomRepository.CreateRoom(&room); err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.Connections[room.ID] = make(map[string]*Connection)
	r.WaitingRoom[room.ID] = make(map[string]*Connection)
	r.mu.Unlock()

	return &room.ID, nil
}

func (r *RoomService) JoinRoom(roomID, userID string) (string, error) {
	exists, err := r.RoomRepository.DoesRoomExist(roomID)
	if err != nil || !*exists {
		return "", fmt.Errorf("room does not exist")
	}

	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return "", err
	}

	// Host is automatically admitted
	if room.HostID == userID {
		return models.ParticipantStatusAdmitted, nil
	}

	return models.ParticipantStatusWaiting, nil
}

func (r *RoomService) HandleWebSocket(roomID, userID string, conn *websocket.Conn) error {
	r.mu.Lock()
	if r.Connections[roomID] == nil {
		r.Connections[roomID] = make(map[string]*Connection)
	}

	connection := &Connection{
		Conn:     conn,
		UserID:   userID,
		JoinedAt: time.Now(),
		Status:   models.ParticipantStatusAdmitted,
	}
	r.Connections[roomID][userID] = connection
	r.mu.Unlock()

	// Notify others about new peer
	r.broadcastToRoom(roomID, models.WebSocketMessage{
		Type: models.WSMessageTypeJoin,
		Payload: map[string]string{
			"userId": userID,
			"status": models.ParticipantStatusAdmitted,
		},
	}, userID)

	defer func() {
		r.removeConnection(roomID, userID)
	}()

	return r.handleMessages(roomID, userID, conn)
}

func (r *RoomService) HandleWaitingRoom(roomID, userID string, conn *websocket.Conn) error {
	r.mu.Lock()
	if r.WaitingRoom[roomID] == nil {
		r.WaitingRoom[roomID] = make(map[string]*Connection)
	}

	connection := &Connection{
		Conn:     conn,
		UserID:   userID,
		JoinedAt: time.Now(),
		Status:   models.ParticipantStatusWaiting,
	}
	r.WaitingRoom[roomID][userID] = connection
	r.mu.Unlock()

	// Notify host about waiting participant
	r.notifyHost(roomID, models.WebSocketMessage{
		Type: "waiting-participant",
		Payload: map[string]string{
			"userId": userID,
		},
	})

	defer func() {
		r.removeFromWaitingRoom(roomID, userID)
	}()

	// Wait for admission decision
	for {
		var msg models.WebSocketMessage
		if err := conn.ReadJSON(&msg); err != nil {
			return err
		}
	}
}

func (r *RoomService) AdmitParticipant(roomID, participantID, hostID string) error {
	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room.HostID != hostID {
		return errors.New("only host can admit participants")
	}

	r.mu.Lock()
	participant, exists := r.WaitingRoom[roomID][participantID]
	if !exists {
		r.mu.Unlock()
		return errors.New("participant not found in waiting room")
	}

	// Move from waiting room to admitted participants
	delete(r.WaitingRoom[roomID], participantID)
	r.Connections[roomID][participantID] = participant
	participant.Status = models.ParticipantStatusAdmitted
	r.mu.Unlock()

	// Notify participant about admission
	msg := models.WebSocketMessage{
		Type:    models.WSMessageTypeAdmitted,
		Payload: map[string]string{"status": "admitted"},
	}
	return participant.Conn.WriteJSON(msg)
}

func (r *RoomService) DenyParticipant(roomID, participantID, hostID string) error {
	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return err
	}

	if room.HostID != hostID {
		return errors.New("only host can deny participants")
	}

	r.mu.Lock()
	participant, exists := r.WaitingRoom[roomID][participantID]
	if !exists {
		r.mu.Unlock()
		return errors.New("participant not found in waiting room")
	}

	// Remove from waiting room
	delete(r.WaitingRoom[roomID], participantID)
	r.mu.Unlock()

	// Notify participant about denial
	msg := models.WebSocketMessage{
		Type:    models.WSMessageTypeDenied,
		Payload: map[string]string{"status": "denied"},
	}
	return participant.Conn.WriteJSON(msg)
}

// Additional helper methods...
func (r *RoomService) GetRoomStatus(roomID, userID string) (*models.RoomStatus, error) {
	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	status := &models.RoomStatus{
		RoomID:       room.ID,
		Name:         room.Name,
		HostID:       room.HostID,
		IsActive:     room.Status == models.RoomStatusActive,
		Participants: len(r.Connections[roomID]),
		WaitingCount: len(r.WaitingRoom[roomID]),
		CreatedAt:    room.CreatedAt,
	}

	return status, nil
}

func (r *RoomService) GetParticipants(roomID, userID string) (*models.Participants, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := &models.Participants{
		Admitted: make([]models.Participant, 0),
		Waiting:  make([]models.Participant, 0),
	}

	// Get admitted participants
	for _, conn := range r.Connections[roomID] {
		result.Admitted = append(result.Admitted, models.Participant{
			UserID:   conn.UserID,
			Username: conn.Username,
			JoinedAt: conn.JoinedAt,
			Status:   models.ParticipantStatusAdmitted,
		})
	}

	// Get waiting participants
	for _, conn := range r.WaitingRoom[roomID] {
		result.Waiting = append(result.Waiting, models.Participant{
			UserID:   conn.UserID,
			Username: conn.Username,
			JoinedAt: conn.JoinedAt,
			Status:   models.ParticipantStatusWaiting,
		})
	}

	return result, nil
}

// broadcastToRoom sends a message to all connections in a room except the sender
func (r *RoomService) broadcastToRoom(roomID string, message interface{}, senderID string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if room exists
	connections, exists := r.Connections[roomID]
	if !exists {
		return fmt.Errorf("room %s not found", roomID)
	}

	// Send to all connections except sender
	for userID, conn := range connections {
		if userID != senderID {
			err := conn.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message to user %s in room %s: %v", userID, roomID, err)
				continue // Continue broadcasting to others even if one fails
			}
		}
	}

	return nil
}

// removeConnection removes a WebSocket connection from a room
func (r *RoomService) removeConnection(roomID string, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if room exists in connections map
	if connections, exists := r.Connections[roomID]; exists {
		// Remove the specific user's connection
		if _, ok := connections[userID]; ok {
			connections[userID].Conn.Close()
			delete(connections, userID)
			log.Printf("Connection removed from room %s. Total connections: %d", roomID, len(connections))
		}

		// Clean up room if empty
		if len(connections) == 0 {
			delete(r.Connections, roomID)
			log.Printf("Room %s removed as it's empty", roomID)
		}
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

// IsUserAdmitted checks if a user is admitted to the room
func (r *RoomService) IsUserAdmitted(roomID, userID string) (bool, error) {
	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return false, err
	}

	// Host is always admitted
	if room.HostID == userID {
		return true, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if user is in admitted connections
	_, isAdmitted := r.Connections[roomID][userID]
	return isAdmitted, nil
}

// LeaveRoom handles a user leaving the room
func (r *RoomService) LeaveRoom(roomID, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user is in admitted connections
	if conn, exists := r.Connections[roomID][userID]; exists {
		// Notify others about the user leaving
		r.broadcastToRoom(roomID, models.WebSocketMessage{
			Type: models.WSMessageTypeLeave,
			Payload: map[string]string{
				"userId": userID,
			},
		}, userID)

		// Close connection and remove from admitted list
		conn.Conn.Close()
		delete(r.Connections[roomID], userID)

		// If room is empty, clean up
		if len(r.Connections[roomID]) == 0 {
			delete(r.Connections, roomID)
			delete(r.WaitingRoom, roomID)
		}

		return nil
	}

	// Check if user is in waiting room
	if conn, exists := r.WaitingRoom[roomID][userID]; exists {
		conn.Conn.Close()
		delete(r.WaitingRoom[roomID], userID)
		return nil
	}

	return errors.New("user not found in room")
}

// removeFromWaitingRoom removes a user from the waiting room
func (r *RoomService) removeFromWaitingRoom(roomID string, userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if waitingRoom, exists := r.WaitingRoom[roomID]; exists {
		if conn, ok := waitingRoom[userID]; ok {
			conn.Conn.Close()
			delete(waitingRoom, userID)
			log.Printf("User %s removed from waiting room %s", userID, roomID)
		}

		// Clean up waiting room if empty
		if len(waitingRoom) == 0 {
			delete(r.WaitingRoom, roomID)
			log.Printf("Waiting room for room %s removed as it's empty", roomID)
		}
	}
}

// notifyHost sends a message to the room host
func (r *RoomService) notifyHost(roomID string, message models.WebSocketMessage) error {
	room, err := r.RoomRepository.GetRoom(roomID)
	if err != nil {
		return err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	// Find host's connection
	if connections, exists := r.Connections[roomID]; exists {
		if hostConn, ok := connections[room.HostID]; ok {
			return hostConn.Conn.WriteJSON(message)
		}
	}

	return errors.New("host not connected")
}

// handleMessages handles incoming WebSocket messages
func (r *RoomService) handleMessages(roomID string, userID string, conn *websocket.Conn) error {
	for {
		var msg models.WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			return err
		}

		switch msg.Type {
		case models.WSMessageTypeOffer,
			models.WSMessageTypeAnswer,
			models.WSMessageTypeIceCandidate:
			// Forward WebRTC signaling messages to other participants
			if err := r.broadcastToRoom(roomID, msg, userID); err != nil {
				log.Printf("error broadcasting message: %v", err)
			}

		case models.WSMessageTypeLeave:
			r.removeConnection(roomID, userID)
			r.broadcastToRoom(roomID, models.WebSocketMessage{
				Type: models.WSMessageTypeLeave,
				Payload: map[string]string{
					"userId": userID,
				},
			}, userID)
			return nil

		default:
			log.Printf("unknown message type: %s", msg.Type)
		}
	}
}

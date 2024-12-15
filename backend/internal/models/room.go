package models

import "time"

type Room struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	HostID    string    `json:"hostId"`
	CreatedAt time.Time `json:"createdAt"`
	Status    int       `json:"status"`
}

type Participant struct {
	UserID   string    `json:"userId"`
	Username string    `json:"username"`
	Name     string    `json:"name"`
	JoinedAt time.Time `json:"joinedAt"`
	Status   string    `json:"status"` // "waiting", "admitted", "denied"
}

type Participants struct {
	Admitted []Participant `json:"admitted"`
	Waiting  []Participant `json:"waiting"`
}

type RoomStatus struct {
	RoomID       string    `json:"roomId"`
	Name         string    `json:"name"`
	HostID       string    `json:"hostId"`
	IsActive     bool      `json:"isActive"`
	Participants int       `json:"participantCount"`
	WaitingCount int       `json:"waitingCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

// WebSocket message types
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// Constants for room status
const (
	RoomStatusActive   = 1
	RoomStatusInactive = 2
)

// Constants for participant status
const (
	ParticipantStatusWaiting  = "waiting"
	ParticipantStatusAdmitted = "admitted"
	ParticipantStatusDenied   = "denied"
)

// Constants for WebSocket message types
const (
	WSMessageTypeJoin              = "join"
	WSMessageTypeLeave             = "leave"
	WSMessageTypeOffer             = "offer"
	WSMessageTypeAnswer            = "answer"
	WSMessageTypeIceCandidate      = "ice-candidate"
	WSMessageTypeAdmitted          = "admitted"
	WSMessageTypeDenied            = "denied"
	WSMessageTypeParticipantUpdate = "participant-update"
)

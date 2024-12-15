package service

import (
	"sync"

	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/session"
)

type AuthService struct {
	AuthRepository *repositories.AuthRepository
	SessionManager *session.SessionManager
}

// RoomService handles room operations and WebRTC signaling
type RoomService struct {
	RoomRepository *repositories.RoomRepository
	mu             sync.RWMutex
	Connections    map[string]map[string]*Connection // roomID -> userID -> Connection
	WaitingRoom    map[string]map[string]*Connection // roomID -> userID -> Connection
}

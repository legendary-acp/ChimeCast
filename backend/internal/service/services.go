package service

import (
	"github.com/gorilla/websocket"
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/session"
)

type AuthService struct {
	AuthRepository *repositories.AuthRepository
	SessionManager *session.SessionManager
}

type RoomService struct {
	RoomRepository *repositories.RoomRepository
	Connections    map[string][]*websocket.Conn
}

package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/legendary-acp/chimecast/internal/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

type RoomHandler struct {
	RoomService *service.RoomService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

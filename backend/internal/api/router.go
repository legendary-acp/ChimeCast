package api

import (
	"github.com/gorilla/mux"
	"github.com/legendary-acp/chimecast/internal/api/handler"
	"github.com/legendary-acp/chimecast/internal/middleware"
	"github.com/legendary-acp/chimecast/internal/service"
	"github.com/legendary-acp/chimecast/internal/session"
)

func NewRouter(authService *service.AuthService, roomService *service.RoomService, sessionManager *session.SessionManager) *mux.Router {
	router := mux.NewRouter()

	authHandler := handler.NewAuthHandler(authService)
	roomHandler := handler.NewRoomHandler(roomService)

	authAPIsV1 := router.PathPrefix("/api/auth/v1").Subrouter()

	authAPIsV1.HandleFunc("/register", authHandler.Register).Methods("POST")
	authAPIsV1.HandleFunc("/login", authHandler.Login).Methods("POST")
	authAPIsV1.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	roomAPIsV1 := router.PathPrefix("/api/room/v1").Subrouter()
	roomAPIsV1.Use(middleware.AuthMiddleware(sessionManager))

	roomAPIsV1.HandleFunc("/", roomHandler.GetAllRooms).Methods("GET")
	roomAPIsV1.HandleFunc("/", roomHandler.CreateRoom).Methods("POST")
	roomAPIsV1.HandleFunc("/{roomID}/join", roomHandler.JoinRoom).Methods("POST")
	roomAPIsV1.HandleFunc("/{roomID}/ws", roomHandler.HandleWebSocket).Methods("GET")

	return router
}

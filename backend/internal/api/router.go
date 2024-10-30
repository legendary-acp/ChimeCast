package api

import (
	"github.com/gorilla/mux"
	"github.com/legendary-acp/chimecast/internal/api/handler"
	"github.com/legendary-acp/chimecast/internal/service"
)

func NewRouter(authService *service.AuthService) *mux.Router {
	router := mux.NewRouter()

	authHandler := handler.NewAuthHandler(authService)

	authAPIsV1 := router.PathPrefix("/api/auth/v1").Subrouter()

	authAPIsV1.HandleFunc("/register", authHandler.Register).Methods("POST")
	authAPIsV1.HandleFunc("/login", authHandler.Login).Methods("POST")
	authAPIsV1.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	return router
}

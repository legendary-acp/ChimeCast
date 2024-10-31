package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/legendary-acp/chimecast/internal/api"
	"github.com/legendary-acp/chimecast/internal/constants"
	"github.com/legendary-acp/chimecast/internal/db"
	"github.com/legendary-acp/chimecast/internal/middleware"
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/service"
	"github.com/legendary-acp/chimecast/internal/session"
)

func main() {
	db, err := db.CreateDB()
	if err != nil {
		log.Fatalln("Unable to Initiate DB")
		return
	}

	sessionManager := session.NewSessionManager()

	authRepository := repositories.NewAuthRepository(db)
	roomRepository := repositories.NewRoomRepositor(db)

	authService := service.NewAuthService(authRepository, sessionManager)
	roomService := service.NewRoomService(roomRepository)

	router := api.NewRouter(authService, roomService, sessionManager)

	handlerWithCors := middleware.CorsMiddleware(router)
	server := &http.Server{
		Addr:    ":" + constants.PORT,
		Handler: handlerWithCors,
	}

	go func() {
		log.Println("Server started on :", constants.PORT)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start Server:", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
}

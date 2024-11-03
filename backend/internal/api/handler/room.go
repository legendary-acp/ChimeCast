package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/service"
	"github.com/legendary-acp/chimecast/internal/utils"
)

func NewRoomHandler(roomService *service.RoomService) *RoomHandler {
	return &RoomHandler{
		RoomService: roomService,
	}
}

func (h *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.RoomService.GetAllRooms()
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, rooms)
	return
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomRequest *models.CreateRoomRequest

	err := json.NewDecoder(r.Body).Decode(&createRoomRequest)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	roomID, err := h.RoomService.CreateRoom(createRoomRequest)
	if err != nil {
		if err.Error() == "name can't be empty" {
			utils.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		} else {
			utils.SendJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"roomID":  *roomID,
		"message": "Room created successfully",
	})
	return
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]

	// Check if the room exists and handle potential errors
	exists, err := h.RoomService.JoinRoom(roomID)
	if err != nil {
		log.Printf("Error joining room %s: %v", roomID, err)
		utils.SendJSONError(w, http.StatusBadRequest, "Could not join the room: "+err.Error())
		return
	}

	// Room existence check
	if !*exists {
		log.Printf("Room ID %s does not exist", roomID)
		utils.WriteJSONResponse(w, http.StatusNotFound, map[string]string{
			"message": "Room ID does not exist",
		})
		return
	}

	// Successful join response
	log.Printf("User joined room %s successfully", roomID)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Joining room",
		"roomID":  roomID,
	})
}

func (h *RoomHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract roomID from URL
	roomID := mux.Vars(r)["roomID"]

	// Upgrade HTTP request to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Printf("Closing WebSocket connection for room %s", roomID)
		conn.Close()
	}()

	// Inform that a user has joined the room
	log.Printf("User joined room: %s", roomID)

	// Call the service to handle the WebSocket connection
	if err := h.RoomService.WebSocketConnection(roomID, conn); err != nil {
		log.Printf("Error handling WebSocket connection for room %s: %v", roomID, err)
		return
	}
}

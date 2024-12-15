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
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var createRoomRequest *models.CreateRoomRequest
	userID := r.Context().Value("userID").(string) // From auth middleware

	err := json.NewDecoder(r.Body).Decode(&createRoomRequest)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	roomID, err := h.RoomService.CreateRoom(createRoomRequest, userID) // Pass hostID
	if err != nil {
		if err.Error() == "name can't be empty" {
			utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		} else {
			utils.SendJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"roomID":  *roomID,
		"message": "Room created successfully",
		"role":    "host",
	})
}

func (h *RoomHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	userID := r.Context().Value("userID").(string)

	status, err := h.RoomService.JoinRoom(roomID, userID)
	if err != nil {
		log.Printf("Error joining room %s: %v", roomID, err)
		utils.SendJSONError(w, http.StatusBadRequest, "Could not join the room: "+err.Error())
		return
	}

	response := map[string]string{
		"roomID":  roomID,
		"status":  status, // "waiting" or "admitted"
		"message": "Joined room successfully",
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *RoomHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	userID := r.Context().Value("userID").(string)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Check if user is admitted to the room
	isAdmitted, err := h.RoomService.IsUserAdmitted(roomID, userID)
	if err != nil {
		log.Printf("Error checking user admission: %v", err)
		return
	}

	if !isAdmitted {
		// If not admitted, put in waiting room queue
		if err := h.RoomService.HandleWaitingRoom(roomID, userID, conn); err != nil {
			log.Printf("Error handling waiting room: %v", err)
		}
		return
	}

	// Handle admitted user's WebSocket connection
	if err := h.RoomService.HandleWebSocket(roomID, userID, conn); err != nil {
		log.Printf("Error handling WebSocket: %v", err)
	}
}

func (h *RoomHandler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	userID := r.Context().Value("userID").(string)

	participants, err := h.RoomService.GetParticipants(roomID, userID)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, participants)
}

func (h *RoomHandler) AdmitParticipant(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	participantID := mux.Vars(r)["userID"]
	hostID := r.Context().Value("userID").(string)

	err := h.RoomService.AdmitParticipant(roomID, participantID, hostID)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Participant admitted successfully",
	})
}

func (h *RoomHandler) LeaveRoom(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	userID := r.Context().Value("userID").(string)

	err := h.RoomService.LeaveRoom(roomID, userID)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Left room successfully",
	})
}

func (h *RoomHandler) DenyParticipant(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	participantID := mux.Vars(r)["userID"]
	hostID := r.Context().Value("userID").(string)

	err := h.RoomService.DenyParticipant(roomID, participantID, hostID)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Participant denied successfully",
	})
}

func (h *RoomHandler) GetRoomStatus(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["roomID"]
	userID := r.Context().Value("userID").(string)

	status, err := h.RoomService.GetRoomStatus(roomID, userID)
	if err != nil {
		utils.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, status)
}

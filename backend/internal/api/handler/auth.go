package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/service"
	"github.com/legendary-acp/chimecast/internal/utils"
)

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var userRegisterRequest models.RegisterRequest

	// Step 1: Parse and decode the request body
	err := json.NewDecoder(r.Body).Decode(&userRegisterRequest)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error decoding request body: %v", err)

		// Respond with a generic error message
		response := map[string]string{
			"error":   "invalid_request",
			"message": "The request body is invalid.",
		}
		utils.WriteJSONResponse(w, http.StatusBadRequest, response)
		return
	}
	err = utils.ValidateUserRegistrationRequest(&userRegisterRequest)
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Error decoding request body: %v", err)

		// Respond with a generic error message
		response := map[string]string{
			"error":   "invalid_request",
			"message": "The request body is invalid.",
		}
		utils.WriteJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Step 2: Register the user
	sessionID, err := a.AuthService.RegisterUser(&userRegisterRequest)
	if err != nil {
		// Log the actual error
		log.Printf("Error registering user: %v", err)

		// Send a generic error message to the client
		response := map[string]string{
			"error":   "registration_failed",
			"message": "Unable to register user.",
		}

		// Determine error type and respond with the appropriate status code
		if errors.Is(err, utils.ErrUserAlreadyExists) {
			utils.WriteJSONResponse(w, http.StatusConflict, response)
		} else {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, response)
		}
		return
	}

	// Step 3: Success response

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    *sessionID,
		HttpOnly: true,
		Path:     "/",
	})

	response := map[string]string{
		"message": "User registered successfully",
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Login and get session ID
	sessionID, err := a.AuthService.Login(&loginRequest)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusUnauthorized, map[string]string{"message": "Invalid Credentials"})
		return
	}

	// Set session ID as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Path:     "/",
	})

	response := map[string]string{
		"message": "Login successful",
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "no session found"})
		return
	}

	// Delete session
	a.AuthService.Logout(cookie.Value)

	// Clear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	response := map[string]string{
		"message": "Logout successful",
	}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

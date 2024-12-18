package service

import (
	"errors"
	"time"

	"github.com/legendary-acp/chimecast/internal/models"
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/session"
	"github.com/legendary-acp/chimecast/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func NewAuthService(authRepositories *repositories.AuthRepository, sessionManager *session.SessionManager) *AuthService {
	return &AuthService{
		AuthRepository: authRepositories,
		SessionManager: sessionManager,
	}
}

func (a *AuthService) RegisterUser(request *models.RegisterRequest) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a user model instance with the hashed password
	user := models.User{
		ID:             utils.CreateNewUUID(),
		Name:           request.Name,
		Username:       request.UserName,
		Email:          request.Email,
		HashedPassword: string(hashedPassword),
		CreatedAt:      time.Now(),
	}

	err = a.AuthRepository.RegisterUser(user)
	if err != nil {
		return nil, err
	}
	sessionID, err := a.SessionManager.CreateSession(user.Username, user.ID)
	if err != nil {
		return nil, errors.New("could not create session")
	}

	return &sessionID, nil
}

func (a *AuthService) Login(request *models.LoginRequest) (string, error) {
	user, err := a.AuthRepository.Login(request.UserName)
	if err != nil {
		return "", err
	}

	// Compare the provided password with the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(request.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate session
	sessionID, err := a.SessionManager.CreateSession(user.Username, user.ID)
	if err != nil {
		return "", errors.New("could not create session")
	}

	return sessionID, nil
}

func (a *AuthService) Logout(sessionID string) error {
	a.SessionManager.DeleteSession(sessionID)
	return nil
}

func (a *AuthService) ValidateAuth(sessionID string) error {
	// Get and validate session
	_, err := a.SessionManager.GetSession(sessionID)

	return err
}

package service

import (
	"github.com/legendary-acp/chimecast/internal/repositories"
	"github.com/legendary-acp/chimecast/internal/session"
)

type AuthService struct {
	AuthRepository *repositories.AuthRepository
	SessionManager *session.SessionManager
}

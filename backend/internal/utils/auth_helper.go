package utils

import (
	"errors"

	"github.com/legendary-acp/chimecast/internal/models"
)

func ValidateUserRegistrationRequest(request *models.RegisterRequest) error {
	if request.Name == "" {
		return errors.New("name cannot be empty")
	}
	if request.Email == "" {
		return errors.New("email cannot be empty")
	}
	if request.Password == "" {
		return errors.New("password cannot be empty")
	}
	if request.UserId == "" {
		return errors.New("userId cannot be empty")
	}
	return nil
}

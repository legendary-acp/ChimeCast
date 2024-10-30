package models

import "time"

type User struct {
	ID             string    `json:"-"`
	UserName       string    `json:"userName"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedPassword"`
	CreatedAt      time.Time `json:"CreatedAt"`
}

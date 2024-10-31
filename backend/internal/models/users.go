package models

import "time"

type User struct {
	ID             string    `json:"-"`
	Username       string    `json:"Username"`
	Name           string    `json:"Name"`
	Email          string    `json:"Email"`
	HashedPassword string    `json:"HashedPassword"`
	CreatedAt      time.Time `json:"CreatedAt"`
}

package models

import (
	"github.com/google/uuid"
)

type Profile struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	Admin bool `json:"admin"`
}

// Stores unique user ID and Name
type UniqueInfo struct {
	ID   uuid.UUID
	Name string
}

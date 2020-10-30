package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// SignedUser properties of an authenticated user
type SignedUser struct {
	UserID   string   `json:"userId"`
	UserType UserType `json:"userType"`
	Token    string   `json:"token"`
}

// User entity definition. Describes any of
// admin, agent, merchant or subscriber
type User struct {
	UserID   uuid.UUID
	UserType UserType

	gorm.Model // embed created_at, deleted_at, updated_at
}

package models

import (
	"time"

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

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

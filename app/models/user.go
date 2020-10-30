package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// SignedUser properties of an authenticated user
type SignedUser struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

// User entity definition
type User struct {
	gorm.Model // embed created_at, deleted_at, updated_at

	ID          uuid.UUID `gorm:"primarykey"`
	FirstName   string
	LastName    string
	Email       string `gorm:"not null;unique"`
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

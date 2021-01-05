package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Subscriber
type Subscriber struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Subscriber) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

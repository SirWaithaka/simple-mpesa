package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Merchant
type Merchant struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	// a merchant is usually assigned a till number they use to accept
	// payments from other customers
	TillNumber string `gorm:"column:till_number;unique"`

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Merchant) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Agent
type Agent struct {
	ID    uuid.UUID
	Email string `gorm:"not null;unique"` // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string `gorm:"not null;unique"`
	PassportNo  string
	Password    string `gorm:"not null"`

	// an agent is usually assigned an agent number that they use for
	// transactions with other customers
	AgentNumber string `gorm:"column:agent_number;unique"`

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Agent) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

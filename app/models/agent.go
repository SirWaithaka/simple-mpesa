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
	// AgentNumber string `gorm:"column:agent_number;unique"`

	// an extra column/property that tells us if the agent is a super agent
	SuperAgent SuperAgentStatus `gorm:"default:'0'"`// PS: bool values dont work well with gorm during updates

	gorm.Model
}

// BeforeCreate hook will be used to add uuid to entity before adding to db
func (u *Agent) BeforeCreate(tx *gorm.DB) error {
	u.ID, _ = uuid.NewV4()
	return nil
}

func (u Agent) IsSuperAgent() bool {
	return u.SuperAgent == IsSuperAgent
}

// SuperAgentStatus
type SuperAgentStatus string

const (
	IsNotSuperAgent = SuperAgentStatus('0')
	IsSuperAgent    = SuperAgentStatus('1')
)

// Not returns the opposite, if that makes sense
func (status SuperAgentStatus) Not() SuperAgentStatus {
	if status == IsNotSuperAgent {
		return IsSuperAgent
	} else {
		return IsNotSuperAgent
	}
}

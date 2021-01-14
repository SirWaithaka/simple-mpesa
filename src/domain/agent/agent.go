package agent

import (
	"github.com/gofrs/uuid"
)

// Agent
type Agent struct {
	ID    uuid.UUID
	Email string // email is used as account number

	FirstName   string
	LastName    string
	PhoneNumber string
	PassportNo  string
	Password    string

	// an extra column/property that tells us if the agent is a super agent
	SuperAgent SuperAgentStatus // PS: bool values dont work well with gorm during updates
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

type Repository interface {
	Add(Agent) (Agent, error)
	Delete(Agent) error
	FetchAll() ([]Agent, error)
	FindByID(uuid.UUID) (Agent, error)
	FindByEmail(string) (Agent, error)
	Update(Agent) error
}

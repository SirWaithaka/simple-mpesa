package repositories

import (
	"simple-mpesa/src/agent"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewAgentRepository(database *storage.Database) *AgentRepository {
	return &AgentRepository{db: database}
}

type AgentRepository struct {
	db *storage.Database
}

func (r AgentRepository) searchBy(row agent.Agent) (agent.Agent, error) {
	var a agent.Agent
	result := r.db.Where(row).First(&a)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return agent.Agent{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return agent.Agent{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return a, nil
}

// Add an agent if not in db.
func (r AgentRepository) Add(a agent.Agent) (agent.Agent, error) {
	// add new agent to agents table, if query return violation of unique key column,
	// we know that the agent with given record exists and return that agent instead
	result := r.db.Model(Agent{}).Create(&a)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return a, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return agent.Agent{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return a, nil
}

// Delete a agent
func (r AgentRepository) Delete(agent agent.Agent) error {
	result := r.db.Delete(&agent)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all agents in db
func (r AgentRepository) FetchAll() ([]agent.Agent, error) {
	var agents []agent.Agent
	result := r.db.Find(&agents)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	// we might not need to return this error
	if result.RowsAffected == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND}
	}

	return agents, nil
}

// FindByID searches agent by primary id
func (r AgentRepository) FindByID(id uuid.UUID) (agent.Agent, error) {
	a, err := r.searchBy(agent.Agent{ID: id})
	return a, err
}

// FindByEmail searches agent by email
func (r AgentRepository) FindByEmail(email string) (agent.Agent, error) {
	a, err := r.searchBy(agent.Agent{Email: email})
	return a, err
}

// Update
func (r AgentRepository) Update(a agent.Agent) error {
	result := r.db.Model(&agent.Agent{}).Where(agent.Agent{ID: a.ID}).Omit("id").Updates(a)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}


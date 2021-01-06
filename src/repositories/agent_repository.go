package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewAgentRepository(database *storage.Database) *Agent {
	return &Agent{db: database}
}

type Agent struct {
	db *storage.Database
}

func (r Agent) searchBy(row models.Agent) (models.Agent, error) {
	var agent models.Agent
	result := r.db.Where(row).First(&agent)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.Agent{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return models.Agent{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return agent, nil
}

// Add an agent if not in db.
func (r Agent) Add(agent models.Agent) (models.Agent, error) {
	// add new agent to agents table, if query return violation of unique key column,
	// we know that the agent with given record exists and return that agent instead
	result := r.db.Model(models.Agent{}).Create(&agent)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return agent, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrUserExists}
		}
		return models.Agent{}, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return agent, nil
}

// Delete a agent
func (r Agent) Delete(agent models.Agent) error {
	result := r.db.Delete(&agent)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// FetchAll gets all agents in db
func (r Agent) FetchAll() ([]models.Agent, error) {
	var agents []models.Agent
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
func (r Agent) FindByID(id uuid.UUID) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{ID: id})
	return agent, err
}

// FindByEmail searches agent by email
func (r Agent) FindByEmail(email string) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{Email: email})
	return agent, err
}

// Update
func (r Agent) Update(agent models.Agent) error {
	var u models.Agent
	result := r.db.Debug().Model(&u).Where(models.Agent{ID: agent.ID}).Omit("id").Updates(agent)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}


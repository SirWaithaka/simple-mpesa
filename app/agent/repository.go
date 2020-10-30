package agent

import (
	"simple-wallet/app/errors"
	"simple-wallet/app/models"
	"simple-wallet/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(models.Agent) (models.Agent, error)
	Delete(models.Agent) error
	GetByID(uuid.UUID) (models.Agent, error)
	GetByEmail(string) (models.Agent, error)
	Update(models.Agent) error
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) searchBy(row models.Agent) (models.Agent, error) {
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
func (r repository) Add(agent models.Agent) (models.Agent, error) {
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
func (r repository) Delete(agent models.Agent) error {
	result := r.db.Delete(&agent)
	if result.Error != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

// GetByID searches agent by primary id
func (r repository) GetByID(id uuid.UUID) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{ID: id})
	return agent, err
}

// GetByEmail searches agent by email
func (r repository) GetByEmail(email string) (models.Agent, error) {
	agent, err := r.searchBy(models.Agent{Email: email})
	return agent, err
}

// Update
func (r repository) Update(agent models.Agent) error {
	var u models.Agent
	result := r.db.Model(&u).Omit("id").Updates(agent)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

package agent

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(models.Agent) (models.Agent, error)
	Delete(models.Agent) error
	FetchAll() ([]models.Agent, error)
	FindByID(uuid.UUID) (models.Agent, error)
	FindByEmail(string) (models.Agent, error)
	Update(models.Agent) error
}

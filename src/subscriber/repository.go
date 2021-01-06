package subscriber

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(models.Subscriber) (models.Subscriber, error)
	Delete(models.Subscriber) error
	FetchAll() ([]models.Subscriber, error)
	FindByID(uuid.UUID) (models.Subscriber, error)
	FindByEmail(string) (models.Subscriber, error)
	Update(models.Subscriber) error
}

package merchant

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(models.Merchant) (models.Merchant, error)
	Delete(models.Merchant) error
	FetchAll() ([]models.Merchant, error)
	FindByID(uuid.UUID) (models.Merchant, error)
	FindByEmail(string) (models.Merchant, error)
	Update(models.Merchant) error
}

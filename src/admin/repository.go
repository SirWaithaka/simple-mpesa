package admin

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(models.Admin) (models.Admin, error)
	Delete(models.Admin) error
	GetByID(uuid.UUID) (models.Admin, error)
	GetByEmail(string) (models.Admin, error)
	Update(models.Admin) error
}

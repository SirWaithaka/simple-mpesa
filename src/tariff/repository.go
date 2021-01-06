package tariff

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(Charge) (Charge, error)
	FetchAll() ([]Charge, error)
	FindByID(uuid.UUID) (Charge, error)
	Get(operation models.TxnOperation, src models.UserType, dest models.UserType) (Charge, error)
	Update(Charge) error
}

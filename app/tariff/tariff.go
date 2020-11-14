package tariff

import (
	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Tariff struct {
	ID uuid.UUID

	Transaction         models.TxnOperation
	SourceUserType      models.UserType
	DestinationUserType models.UserType
	Charge              models.Cents

	gorm.Model
}

func (t *Tariff) BeforeCreate(tx *gorm.DB) error {
	t.ID, _ = uuid.NewV4()
	return nil
}

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

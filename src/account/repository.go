package account

import (
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

type Repository interface {
	GetAccountByUserID(uuid.UUID) (models.Account, error)
	UpdateBalance(amount models.Cents, userID uuid.UUID) (models.Account, error)

	Create(userId uuid.UUID) (models.Account, error)
}

package transaction

import (
	"simple-mpesa/src/models"
)

type Repository interface {
	Add(models.Transaction) (models.Transaction, error)
}

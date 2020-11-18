package account_test

import (
	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type TestAccountant struct {
	DebitAccountFunc  func(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error)
	CreditAccountFunc func(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error)
}

func (ta TestAccountant) DebitAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
	return ta.DebitAccountFunc(userID, amount, reason)
}

func (ta TestAccountant) CreditAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
	return ta.CreditAccountFunc(userID, amount, reason)
}

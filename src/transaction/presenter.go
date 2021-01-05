package transaction

import (
	"time"

	"simple-mpesa/src/data"
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

func parseToTransaction(newTx data.TransactionContract) *models.Transaction {
	id, _ := uuid.NewV4()

	return &models.Transaction{
		ID:        id,
		Operation: newTx.TxnOperation,
		Timestamp: time.Now(),
		Amount:    newTx.Amount,
		UserID:    newTx.UserID,
		AccountID: newTx.AccountID,
	}
}

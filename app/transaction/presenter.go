package transaction

import (
	"time"

	"simple-wallet/app/data"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

func parseToTransaction(newTx data.TransactionContract) *models.Transaction {
	id, _ := uuid.NewV4()

	return &models.Transaction{
		ID:        id,
		Type:      newTx.TxType,
		Timestamp: time.Now(),
		Amount:    newTx.Amount,
		UserID:    newTx.UserID,
		AccountID: newTx.AccountID,
	}
}

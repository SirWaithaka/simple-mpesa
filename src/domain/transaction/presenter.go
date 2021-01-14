package transaction

import (
	"time"

	"simple-mpesa/src/data"

	"github.com/gofrs/uuid"
)

func parseToTransaction(newTx data.TransactionContract) *Statement {
	id, _ := uuid.NewV4()

	return &Statement{
		ID:        id,
		Operation: newTx.TxnOperation,
		Timestamp: time.Now(),
		Amount:    newTx.Amount,
		UserID:    newTx.UserID,
		AccountID: newTx.AccountID,
	}
}

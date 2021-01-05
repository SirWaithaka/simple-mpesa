package account

import (
	"time"

	"simple-mpesa/src/data"
	"simple-mpesa/src/models"

	"github.com/gofrs/uuid"
)

func parseTransactionDetails(userId uuid.UUID, acc models.Account, txnOp models.TxnOperation, timestamp time.Time) *data.TransactionContract {
	return &data.TransactionContract{
		UserID:       userId,
		AccountID:    acc.ID,
		Amount:       acc.Balance(),
		TxnOperation: txnOp,
		Timestamp:    timestamp,
	}
}

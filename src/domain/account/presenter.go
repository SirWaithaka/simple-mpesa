package account

import (
	"time"

	"simple-mpesa/src/data"
	"simple-mpesa/src/value_objects"

	"github.com/gofrs/uuid"
)

func parseTransactionDetails(userId uuid.UUID, acc Account, txnOp value_objects.TxnOperation, timestamp time.Time) *data.TransactionContract {
	return &data.TransactionContract{
		UserID:       userId,
		AccountID:    acc.ID,
		Amount:       acc.Balance(),
		TxnOperation: txnOp,
		Timestamp:    timestamp,
	}
}

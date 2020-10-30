package account

import (
	"time"

	"simple-wallet/app/data"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

func parseTransactionDetails(userId uuid.UUID, acc models.Account, txType string, timestamp time.Time) *data.TransactionContract {
	return &data.TransactionContract{
		UserID:    userId,
		AccountID: acc.ID,
		Amount:    acc.Balance(),
		TxType:    txType,
		Timestamp: timestamp,
	}
}

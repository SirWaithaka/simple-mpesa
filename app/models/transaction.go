package models

import (
	"time"

	"github.com/gofrs/uuid"
)

const (
	TxTypeDeposit    = "deposit"
	TxTypeWithdrawal = "withdrawal"
	TxTypeBalance    = "balance_enquiry"
)

type Transaction struct {
	ID        uuid.UUID `json:"transactionId"`
	Type      string    `json:"transactionType"`
	Timestamp time.Time `json:"timestamp"`
	Amount    float64   `json:"amount"`
	UserID    uuid.UUID `json:"userId"`
	AccountID uuid.UUID `json:"accountId"`
}

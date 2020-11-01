package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type TxType string

const (
	TxTypeDeposit    = TxType("deposit")
	TxTypeWithdrawal = TxType("withdrawal")
	TxTypeBalance    = TxType("balance_enquiry")
)

type Transaction struct {
	ID        uuid.UUID
	Type      string
	Timestamp time.Time
	Amount    float64
	UserID    uuid.UUID
	AccountID uuid.UUID
}

package account

import (
	"time"

	"simple-mpesa/src/domain/value_objects"

	"github.com/gofrs/uuid"
)

type TransactionType string

const (
	TxnTypeCredit = TransactionType("CREDIT")
	TxnTypeDebit  = TransactionType("DEBIT")
)

type Statement struct {
	ID           uuid.UUID
	Operation    value_objects.TxnOperation
	DebitAmount  float64
	CreditAmount float64
	UserID       uuid.UUID
	AccountID    uuid.UUID
	CreatedAt    time.Time
}

type StatementRepository interface {
	Add(Statement) (Statement, error)
	GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error)
}

package statement

import (
	"time"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Type string

const (
	TypeCredit = Type("CREDIT")
	TypeDebit  = Type("DEBIT")
)

type Statement struct {
	ID           uuid.UUID
	Operation    models.TxnOperation
	DebitAmount  float64
	CreditAmount float64
	UserID       uuid.UUID
	AccountID    uuid.UUID
	CreatedAt    time.Time
}

func (s *Statement) BeforeCreate(tx *gorm.DB) error {
	s.ID, _ = uuid.NewV4()
	return nil
}

func (Statement) TableName() string {
	return "statements"
}

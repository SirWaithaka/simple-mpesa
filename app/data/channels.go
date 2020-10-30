package data

import (
	"time"

	"github.com/gofrs/uuid"
)

// UserContract describe the characteristics of data that should
// be passed along in channels for when a user is created or something.
type UserContract struct {
	UserID uuid.UUID
}

type ChanNewUsers struct {
	Channel chan UserContract
	Reader  <-chan UserContract
	Writer  chan<- UserContract
}

// TransactionContract represents the type of data
// required to record a new transaction in the database.
type TransactionContract struct {
	UserID    uuid.UUID
	AccountID uuid.UUID
	Amount    float64
	TxType    string // transaction type
	Timestamp time.Time
}

type ChanNewTransactions struct {
	Channel chan TransactionContract
	Reader  <-chan TransactionContract
	Writer  chan<- TransactionContract
}

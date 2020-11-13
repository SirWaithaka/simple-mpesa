package data

import (
	"time"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

// CustomerContract describe the characteristics of data that should
// be passed along in channels for when a user is created or something.
type CustomerContract struct {
	UserID uuid.UUID
}

type ChanNewCustomers struct {
	Channel chan CustomerContract
	Reader  <-chan CustomerContract
	Writer  chan<- CustomerContract
}

// TransactionContract represents the type of data
// required to record a new transaction in the database.
type TransactionContract struct {
	UserID       uuid.UUID
	AccountID    uuid.UUID
	Amount       float64
	TxnOperation models.TxnOperation // transaction operation type
	Timestamp    time.Time
}

type ChanNewTransactions struct {
	Channel chan TransactionContract
	Reader  <-chan TransactionContract
	Writer  chan<- TransactionContract
}

// type ChanNewTxnEvents struct {
// 	Channel chan models.TxnEvent
// 	Reader  <-chan models.TxnEvent
// 	Writer  chan<- models.TxnEvent
// }

package transaction

import (
	"time"

	"simple-mpesa/src/value_objects"

	"github.com/gofrs/uuid"
)

type Statement struct {
	ID        uuid.UUID
	Operation value_objects.TxnOperation
	Timestamp time.Time
	Amount    float64
	UserID    uuid.UUID
	AccountID uuid.UUID
}

// TxnCustomer is a description of a customer involved in a transaction. We can describe them
// by their user id and user type; We have defined a customer being an agent, merchant or subscriber.
type TxnCustomer struct {
	UserID   uuid.UUID
	UserType value_objects.UserType
}

//
type Transaction struct {
	Source      TxnCustomer // where money is coming from
	Destination TxnCustomer // where money is going
	// we can further use this field to describe the specific type of transaction/transfer
	TxnOperation value_objects.TxnOperation
	// amount of money if shillings being transacted
	Amount value_objects.Shillings
}

type Repository interface {
	Add(Statement) (Statement, error)
}

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
	Type      TxType
	Timestamp time.Time
	Amount    float64
	UserID    uuid.UUID
	AccountID uuid.UUID
}

// TxnOperation is a description of a transaction operation. We have defined operations
// as deposit, withdrawal and transfer. In the end all 3 operations can be modelled as one;
// "transfer" operations.
//
// Example:
// 1. During a deposit, money is moved from an agent's account to the depositor's account
// 2. During a withdrawal, money is moved from the customer withdrawing to the agent's account.
//
// A transfer operation/transaction, needs to have a source and destination and the amount being
// transferred.
type TxnOperation struct {
	Source      TxnCustomer // where money is coming from
	Destination TxnCustomer // where money is going
	// we can further use this field to describe the specific type of transaction/transfer
	TxnType TxType
	// amount of money if shillings being transacted
	Amount Shillings
}

// TxnCustomer is a description of a customer involved in a transaction. We can describe them
// by their user id and user type; We have defined a customer being an agent, merchant or subscriber.
type TxnCustomer struct {
	UserID   uuid.UUID
	UserType UserType
}

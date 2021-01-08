package account

import (
	"simple-mpesa/src/value_objects"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// Status (active,dormant,frozen,suspended)
type Status string

// AccountType (savings,current,utility)
type Type string

const (
	StatusActive    = Status("active")
	StatusDormant   = Status("dormant")
	StatusFrozen    = Status("frozen")
	StatusSuspended = Status("suspended")
)

const (
	// different types of accounts a user could hold
	// we will use current account only.
	AccTypeSavings = Type("savings")
	AccTypeCurrent = Type("current")
	AccTypeUtility = Type("utility")
)

// Account entity definition
type Account struct {
	ID uuid.UUID

	// balance will be stored in cents
	AvailableBalance value_objects.Cents

	Status      Status
	AccountType Type
	UserID      uuid.UUID // a user can only have one account

	gorm.Model
}

// Balance converts balance from cents
func (acc Account) Balance() float64 {
	return float64(acc.AvailableBalance / 100)
}

// Credit add an amount to account balance and return it
func (acc Account) Credit(amount value_objects.Cents) value_objects.Cents {
	// convert incoming amount into cents and add to account balance
	return amount + acc.AvailableBalance
}

// Debit subtract an amount from account balance and return it
func (acc Account) Debit(amount value_objects.Cents) value_objects.Cents {
	// convert incoming amount into cents and subtract to account balance
	return acc.AvailableBalance - amount
}

// IsBalanceLessThanAmount converts amount into cents and returns true if balance is less than amount
func (acc Account) IsBalanceLessThanAmount(amount value_objects.Cents) bool {
	return acc.AvailableBalance < amount
}

// Repository describes crud operation on the account
type Repository interface {
	GetAccountByUserID(uuid.UUID) (Account, error)
	UpdateBalance(amount value_objects.Cents, userID uuid.UUID) (Account, error)

	Create(userId uuid.UUID) (Account, error)
}

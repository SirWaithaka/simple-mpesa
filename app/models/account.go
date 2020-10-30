package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// AccountStatus (active,dormant,frozen,suspended)
type AccountStatus string

// AccountType (savings,current,utility)
type AccountType string

const (
	StatusActive    = AccountStatus("active")
	StatusDormant   = AccountStatus("dormant")
	StatusFrozen    = AccountStatus("frozen")
	StatusSuspended = AccountStatus("suspended")
)

const (
	// different types of accounts a user could hold
	// we will use current account only.
	AccTypeSavings = AccountType("savings")
	AccTypeCurrent = AccountType("current")
	AccTypeUtility = AccountType("utility")
)

// Account entity definition
type Account struct {
	ID uuid.UUID

	// balance will be stored in cents
	AvailableBalance uint `gorm:"column:available_balance"`

	Status      AccountStatus `gorm:"column:status"`
	AccountType AccountType   `gorm:"column:account_type"`
	UserID      uuid.UUID     `gorm:"column:user_id;not null;unique"` // a user can only have one account

	gorm.Model
}

// Balance converts balance from cents
func (acc Account) Balance() float64 {
	return float64(acc.AvailableBalance / 100)
}

// Credit add an amount to account balance and return it
func (acc Account) Credit(amount uint) uint {
	// convert incoming amount into cents and add to account balance
	return (amount * 100) + acc.AvailableBalance
}

// Debit subtract an amount from account balance and return it
func (acc Account) Debit(amount uint) uint {
	// convert incoming amount into cents and subtract to account balance
	return acc.AvailableBalance - (amount * 100)
}

// IsBalanceLessThanAmount converts amount into cents and returns true if balance is less than amount
func (acc Account) IsBalanceLessThanAmount(amount uint) bool {
	return acc.AvailableBalance < (amount * 100)
}

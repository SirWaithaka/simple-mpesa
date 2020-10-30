package errors

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
)

const (
	AccountNotCreated          = ErrorMessage("user's account has not been created, report issue")
	DepositAmountBelowMinimum  = ErrorMessage("cannot deposit amounts less than")
	WithdrawAmountBelowMinimum = ErrorMessage("cannot withdraw amounts less than")
	WithdrawAmountAboveBalance = ErrorMessage("cannot withdraw amount, account balance not enough")
)

// ErrUserHasAccount
func ErrUserHasAccount(userID, accountID uuid.UUID) ErrorMessage {
	return ErrorMessage(fmt.Sprintf("user %v has account with id %v", userID, accountID))
}

// ErrAccountAccess ...
type ErrAccountAccess struct {
	Reason  string
	message string
}

func (err ErrAccountAccess) Error() string {
	msg := fmt.Sprintf("couldn't access account. %v", err.Reason)
	return msg
}

// errAmountBelowMinimum
type errAmountBelowMinimum struct {
	MinAmount uint // minimum amount allowable for deposit or withdraw
	Message   ErrorMessage
}

func (err errAmountBelowMinimum) Error() string {
	return string(err.Message) + " " + strconv.Itoa(int(err.MinAmount))
}

func ErrAmountBelowMinimum(min uint, message ErrorMessage) error {
	return errAmountBelowMinimum{MinAmount: min, Message: message}
}

// ErrNotEnoughBalance
type ErrNotEnoughBalance struct {
	Message ErrorMessage
	Amount  uint
	Balance float64
}

func (err ErrNotEnoughBalance) Error() string {
	return fmt.Sprintf("%s. Amount: %v Balance: %v", string(err.Message), err.Amount, err.Balance)
}

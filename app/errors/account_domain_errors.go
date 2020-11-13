package errors

import (
	"fmt"
	"strconv"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

const (
	AccountNotCreated          = ERMessage("user's account has not been created, report issue")
	DepositAmountBelowMinimum  = ERMessage("cannot deposit amounts less than")
	WithdrawAmountBelowMinimum = ERMessage("cannot withdraw amounts less than")
	TransferAmountBelowMinimum = ERMessage("cannot transfer amounts less than")
	DebitAmountAboveBalance    = ERMessage("cannot debit amount, account balance not enough")

	UserCantHaveAccount = ERMessage("user is not allowed to hold an account")
)

// ErrUserHasAccount
func ErrUserHasAccount(userID, accountID uuid.UUID) ERMessage {
	return ERMessage(fmt.Sprintf("user %v has account with id %v", userID, accountID))
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
	MinAmount models.Shillings // minimum amount allowable for deposit or withdraw
	Message   ERMessage
}

func (err errAmountBelowMinimum) Error() string {
	return string(err.Message) + " " + strconv.Itoa(int(err.MinAmount))
}

func ErrAmountBelowMinimum(min models.Shillings, message ERMessage) error {
	return errAmountBelowMinimum{MinAmount: min, Message: message}
}

// ErrNotEnoughBalance
type ErrNotEnoughBalance struct {
	Message ERMessage
	Amount  models.Shillings
	Balance float64
}

func (err ErrNotEnoughBalance) Error() string {
	return fmt.Sprintf("%s. Amount: %v Balance: %v", string(err.Message), err.Amount, err.Balance)
}

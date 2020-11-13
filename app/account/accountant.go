package account

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type Accountant interface {
	DebitAccount(userID uuid.UUID, amount models.Shillings) (float64, error)
	CreditAccount(userID uuid.UUID, amount models.Shillings) (float64, error)
}

func NewAccountant(accountRepo Repository) Accountant {
	return &accountant{accountRepo}
}

type accountant struct {
	repository Repository
}

func (a accountant) isUserAccAccessible(userID uuid.UUID) (*models.Account, error) {
	acc, err := a.repository.GetAccountByUserID(userID)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return nil, errors.Error{Message: errors.AccountNotCreated, Err: err}
	} else if err != nil {
		return nil, err
	}

	if acc.Status == models.StatusFrozen || acc.Status == models.StatusSuspended {
		e := errors.ErrAccountAccess{Reason: string(acc.Status)}
		return nil, errors.Error{Err: e}
	}

	return &acc, nil

}

func (a accountant) CreditAccount(userID uuid.UUID, amount models.Shillings) (float64, error) {
	acc, err := a.isUserAccAccessible(userID)
	if err != nil {
		return 0, err
	}

	// update balance with amount: add amount
	amt := acc.Credit(amount)
	*acc, err = a.repository.UpdateBalance(amt, userID)
	if err != nil {
		return 0, err
	}

	// a.postTransactionDetails(userID, *acc, models.TxnOpDeposit)
	return acc.Balance(), nil
}

func (a accountant) DebitAccount(userID uuid.UUID, amount models.Shillings) (float64, error) {
	acc, err := a.isUserAccAccessible(userID)
	if err != nil {
		return 0, err
	}

	// check that balance is more than amount
	if acc.IsBalanceLessThanAmount(amount) {
		e := errors.ErrNotEnoughBalance{
			Message: errors.DebitAmountAboveBalance,
			Amount:  amount,
			Balance: acc.Balance(),
		}
		return 0, errors.Error{Err: e}
	}

	// update balance with amount: subtract amount
	amt := acc.Debit(amount)
	*acc, err = a.repository.UpdateBalance(amt, userID)
	if err != nil {
		return 0, err
	}

	// a.postTransactionDetails(userID, *acc, models.TxnOpWithdrawal)
	return acc.Balance(), nil
}

package account

import (
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"

	"github.com/gofrs/uuid"
)

type Accountant interface {
	DebitAccount(userID uuid.UUID, amount value_objects.Cents, reason value_objects.TxnOperation) (float64, error)
	CreditAccount(userID uuid.UUID, amount value_objects.Cents, reason value_objects.TxnOperation) (float64, error)
}

func NewAccountant(accountRepo Repository, ledger Ledger) Accountant {
	return &accountant{repository: accountRepo, ledger: ledger}
}

type accountant struct {
	ledger     Ledger
	repository Repository
}

func (a accountant) isUserAccAccessible(userID uuid.UUID) (*Account, error) {
	acc, err := a.repository.GetAccountByUserID(userID)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return nil, errors.Error{Message: errors.AccountNotCreated, Err: err}
	} else if err != nil {
		return nil, err
	}

	if acc.Status == StatusFrozen || acc.Status == StatusSuspended {
		e := errors.ErrAccountAccess{Reason: string(acc.Status)}
		return nil, errors.Error{Err: e}
	}

	return &acc, nil

}

func (a accountant) CreditAccount(userID uuid.UUID, amount value_objects.Cents, reason value_objects.TxnOperation) (float64, error) {
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

	err = a.ledger.Record(userID, *acc, reason, amount.ToShillings(), TxnTypeCredit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

func (a accountant) DebitAccount(userID uuid.UUID, amount value_objects.Cents, reason value_objects.TxnOperation) (float64, error) {
	acc, err := a.isUserAccAccessible(userID)
	if err != nil {
		return 0, err
	}

	// check that balance is more than amount
	if acc.IsBalanceLessThanAmount(amount) {
		e := errors.ErrNotEnoughBalance{
			Message: errors.DebitAmountAboveBalance,
			Amount:  amount.ToShillings(),
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

	err = a.ledger.Record(userID, *acc, reason, amount.ToShillings(), TxnTypeDebit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

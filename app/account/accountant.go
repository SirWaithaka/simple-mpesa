package account

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/statement"

	"github.com/gofrs/uuid"
)

type Accountant interface {
	DebitAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error)
	CreditAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error)
}

func NewAccountant(accountRepo Repository, ledger statement.Ledger) Accountant {
	return &accountant{repository: accountRepo, ledger: ledger}
}

type accountant struct {
	ledger     statement.Ledger
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

func (a accountant) CreditAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
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

	err = a.ledger.Record(userID, *acc, reason, amount.ToShillings(), statement.TypeCredit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

func (a accountant) DebitAccount(userID uuid.UUID, amount models.Cents, reason models.TxnOperation) (float64, error) {
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

	err = a.ledger.Record(userID, *acc, reason, amount.ToShillings(), statement.TypeDebit)
	if err != nil {
		return 0, err
	}

	return acc.Balance(), nil
}

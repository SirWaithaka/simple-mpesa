package account

import (
	"log"
	"time"

	"simple-wallet/app/data"
	"simple-wallet/app/errors"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

const (
	minimumDepositAmount    = 10 // least possible amount that can be deposited into an account
	minimumWithdrawalAmount = 1  // least possible amount that can be withdrawn from an account
)

type Interactor interface {
	GetBalance(userId uuid.UUID) (float64, error)
	Deposit(userId uuid.UUID, amount uint) (float64, error)
	Withdraw(userId uuid.UUID, amount uint) (float64, error)
}

func NewInteractor(repository Repository, usersChan data.ChanNewUsers, transChan data.ChanNewTransactions) Interactor {
	intr := &interactor{
		repository:          repository,
		usersChannel:        usersChan,
		transactionsChannel: transChan,
	}

	go intr.listenOnNewUsers()

	return intr
}

type interactor struct {
	repository          Repository
	usersChannel        data.ChanNewUsers
	transactionsChannel data.ChanNewTransactions
}

func (i interactor) isUserAccAccessible(userID uuid.UUID) (*models.Account, error) {
	acc, err := i.repository.GetAccountByUserID(userID)
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

// CreateAccount creates an account for a certain user
func (i interactor) CreateAccount(userId uuid.UUID) (models.Account, error) {
	acc, err := i.repository.Create(userId)
	if err != nil {
		return models.Account{}, err
	}
	return acc, nil
}

// GetBalance fetches the user's account balance
func (i interactor) GetBalance(userId uuid.UUID) (float64, error) {
	acc, err := i.isUserAccAccessible(userId)
	if err != nil {
		return 0, err
	}

	i.postTransactionDetails(userId, *acc, models.TxTypeBalance)
	return acc.Balance(), nil
}

// Deposit credits a user's account with an amount
func (i interactor) Deposit(userId uuid.UUID, amount uint) (float64, error) {
	if amount < 10 {
		e := errors.ErrAmountBelowMinimum(minimumDepositAmount, errors.DepositAmountBelowMinimum)
		return 0, errors.Error{Err: e}
	}

	acc, err := i.isUserAccAccessible(userId)
	if err != nil {
		return 0, err
	}

	// update balance with amount: add amount
	amt := acc.Credit(amount)
	*acc, err = i.repository.UpdateBalance(amt, userId)
	if err != nil {
		return 0, err
	}

	i.postTransactionDetails(userId, *acc, models.TxTypeDeposit)
	return acc.Balance(), nil
}

// Withdraw debits a user's account with an amount
func (i interactor) Withdraw(userId uuid.UUID, amount uint) (float64, error) {
	if amount < 10 {
		e := errors.ErrAmountBelowMinimum(minimumWithdrawalAmount, errors.WithdrawAmountBelowMinimum)
		return 0, errors.Error{Err: e}
	}

	acc, err := i.isUserAccAccessible(userId)
	if err != nil {
		return 0, err
	}

	// we can implement a double withdrawal check here. That will prevent a user from
	// withdrawing same amount twice within a stipulated time interval because of system lag.

	// check that balance is more than amount
	if acc.IsBalanceLessThanAmount(amount) {
		e := errors.ErrNotEnoughBalance{
			Message: errors.WithdrawAmountAboveBalance,
			Amount:  amount,
			Balance: acc.Balance(),
		}
		return 0, errors.Error{Err: e}
	}

	// update balance with amount: subtract amount
	amt := acc.Debit(amount)
	log.Printf("new amount %v", amt)
	*acc, err = i.repository.UpdateBalance(amt, userId)
	if err != nil {
		return 0, err
	}

	i.postTransactionDetails(userId, *acc, models.TxTypeWithdrawal)
	return acc.Balance(), nil
}

func (i interactor) postTransactionDetails(userId uuid.UUID, acc models.Account, txType string) {
	timestamp := time.Now()
	newTransaction := parseTransactionDetails(userId, acc, txType, timestamp)

	go func() { i.transactionsChannel.Writer <- *newTransaction }()
}

func (i interactor) listenOnNewUsers() {
	for {
		select {
		case user := <-i.usersChannel.Reader:
			acc, err := i.CreateAccount(user.UserID)
			if err != nil {
				// we need to log this error
				log.Printf("error happened while creating account %v", err)
				return
			}
			// we log the account details if created
			log.Printf("account with id %v has been created successfully for userID %v", acc.ID, user.UserID)
			return
		}
	}
}

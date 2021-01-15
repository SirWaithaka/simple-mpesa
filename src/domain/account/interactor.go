package account

import (
	"log"
	"time"

	"simple-mpesa/src/data"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"

	"github.com/gofrs/uuid"
)

const miniStatementCount = uint(5)

type Interactor interface {
	GetBalance(userID uuid.UUID) (float64, error)
	GetStatement(userId uuid.UUID) ([]Statement, error)
}

func NewInteractor(repository Repository, statementRepo StatementRepository, custChan data.ChanNewCustomers, transChan data.ChanNewTransactions) Interactor {
	intr := &interactor{
		repository:          repository,
		statementRepo:       statementRepo,
		customersChannel:    custChan,
		transactionsChannel: transChan,
	}

	go intr.listenOnNewUsers()

	return intr
}

type interactor struct {
	repository          Repository
	statementRepo       StatementRepository
	customersChannel    data.ChanNewCustomers
	transactionsChannel data.ChanNewTransactions
}

func (i interactor) isUserAccAccessible(userID uuid.UUID) (*Account, error) {
	acc, err := i.repository.GetAccountByUserID(userID)
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

// CreateAccount creates an account for a certain user
func (i interactor) CreateAccount(userId uuid.UUID) (Account, error) {
	acc, err := i.repository.Create(userId)
	if err != nil {
		return Account{}, err
	}
	return acc, nil
}

// GetBalance fetches the user's account balance
func (i interactor) GetBalance(userId uuid.UUID) (float64, error) {
	acc, err := i.isUserAccAccessible(userId)
	if err != nil {
		return 0, err
	}

	// i.postTransactionDetails(userId, *acc, models.TxTypeBalance)
	return acc.Balance(), nil
}

func (i interactor) GetStatement(userID uuid.UUID) ([]Statement, error) {
	now := time.Now()
	transactions, err := i.statementRepo.GetStatements(userID, now, miniStatementCount)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (i interactor) postTransactionDetails(userId uuid.UUID, acc Account, txnOp value_objects.TxnOperation) {
	timestamp := time.Now()
	newTransaction := parseTransactionDetails(userId, acc, txnOp, timestamp)

	go func() { i.transactionsChannel.Writer <- *newTransaction }()
}

func (i interactor) listenOnNewUsers() {
	for {
		select {
		case customer := <-i.customersChannel.Reader:
			acc, err := i.CreateAccount(customer.UserID)
			if err != nil {
				// we need to log this error
				log.Printf("error happened while creating account %v", err)
				continue
			}
			// we log the account details if created
			log.Printf("account with id %v has been created successfully for customerID %v", acc.ID, customer.UserID)
		}
	}
}

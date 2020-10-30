package transaction

import (
	"log"
	"time"

	"simple-wallet/app/data"
	"simple-wallet/app/errors"
	"simple-wallet/app/models"

	"github.com/gofrs/uuid"
)

type Interactor interface {
	AddTransaction(models.Transaction) error
	GetStatement(userId uuid.UUID) (*[]models.Transaction, error)
}

type interactor struct {
	repository   Repository
	transChannel data.ChanNewTransactions
}

func NewInteractor(repository Repository, transChan data.ChanNewTransactions) Interactor {
	intr := &interactor{
		repository:   repository,
		transChannel: transChan,
	}

	go intr.listenOnTransactions()

	return intr
}

func (i interactor) AddTransaction(tx models.Transaction) error {
	_, err := i.repository.Add(tx)
	if err != nil {
		// if we get an error we are going to add the
		// transaction into a buffer object that will
		// retry adding the transaction at a later time

		return err
	}
	return nil
}

func (i interactor) GetStatement(userId uuid.UUID) (*[]models.Transaction, error) {
	now := time.Now()
	transactions, err := i.repository.GetTransactions(userId, now, 5)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (i interactor) listenOnTransactions() {
	for {
		select {
		case tx := <-i.transChannel.Reader:
			transaction := parseToTransaction(tx)

			err := i.AddTransaction(*transaction)
			if err != nil { // if we get an error, it is unexpected, we log it
				log.Printf("error happened when adding transaction to db %v", err.(errors.Error).Err)
				return
			}
			log.Printf("Transaction %v has been successfully added.", transaction.ID)
		}
	}
}

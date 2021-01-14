package transaction

import (
	"log"

	"simple-mpesa/src/data"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
)

const (
	minimumDepositAmount    = value_objects.Shillings(10) // least possible amount that can be deposited into an account
	minimumWithdrawalAmount = value_objects.Shillings(1)  // least possible amount that can be withdrawn from an account
	minimumTransferAmount   = value_objects.Shillings(10) // least possible amount that can be transferred to another account
)

type Interactor interface {
	AddTransaction(Statement) error
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

	go intr.listenOnCreatedTransactions()

	return intr
}

func (i interactor) AddTransaction(tx Statement) error {
	_, err := i.repository.Add(tx)
	if err != nil {
		// if we get an error we are going to add the
		// transaction into a buffer object that will
		// retry adding the transaction at a later time

		return err
	}
	return nil
}

func (i interactor) listenOnCreatedTransactions() {
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

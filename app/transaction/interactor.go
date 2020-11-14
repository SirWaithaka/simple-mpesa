package transaction

import (
	"log"

	"simple-mpesa/app/data"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
)

const (
	minimumDepositAmount    = models.Shillings(10) // least possible amount that can be deposited into an account
	minimumWithdrawalAmount = models.Shillings(1)  // least possible amount that can be withdrawn from an account
	minimumTransferAmount   = models.Shillings(10) // least possible amount that can be transferred to another account
)

type Interactor interface {
	AddTransaction(models.Transaction) error
}

type interactor struct {
	repository       Repository
	transChannel     data.ChanNewTransactions
	// txnEventsChannel data.ChanNewTxnEvents
}

func NewInteractor(repository Repository, transChan data.ChanNewTransactions) Interactor {
	intr := &interactor{
		repository:       repository,
		transChannel:     transChan,
	}

	go intr.listenOnCreatedTransactions()

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

// func (i interactor) listenOnTxnEvents() {
// 	for {
// 		select {
// 		case event := <-i.txnEventsChannel.Reader:
// 			err := i.transact(event)
// 			if err != nil { // if we get an error processing the transaction event, we can log, or send to customer
// 				log.Println(err)
// 				continue
// 			}
// 		}
// 	}
// }

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

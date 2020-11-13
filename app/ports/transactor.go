package ports

import (
	"simple-mpesa/app/customer"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/transaction"

	"github.com/gofrs/uuid"
)

type TransactorPort interface {
	Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error
	Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error
	Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error
}

func NewTransactor(finder customer.Finder, transactor transaction.Transactor) TransactorPort {
	return &transactorAdapter{
		customerFinder: finder,
		transactor:     transactor,
	}
}

type transactorAdapter struct {
	customerFinder customer.Finder
	transactor     transaction.Transactor
}

// Deposit is a transaction between a customer and an agent. The customer's account is credited from the
// agent's account. Money moves from the agent's account to the customer's account.
func (tr transactorAdapter) Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := tr.customerFinder.FindAgentByEmail(agentNumber)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return err
	}


	tx := transaction.Transaction{
		Source: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},
		Destination: depositor,

		TxnOperation: models.TxnOpDeposit,
		Amount:       amount,
	}
	 err = tr.transactor.Transact(tx)
	 if err != nil {
	 	return err
	 }

	return nil
}

// Withdraw is a transaction between a customer and an agent. The customer's account is debited and the
// agent's account credited. Money moves from the customer's account to the agent's account.
func (tr transactorAdapter) Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := tr.customerFinder.FindAgentByEmail(agentNumber)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return err
	}

	tx := transaction.Transaction{
		Source: withdrawer,
		Destination: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},

		TxnOperation: models.TxnOpWithdrawal,
		Amount:       amount,
	}
	err = tr.transactor.Transact(tx)
	if err != nil {
		return err
	}

	return nil
}

// Transfer is a transaction describing a general movement of funds from a customer to another customer. One customer's
// account is debited (the source) and the other customer's account credited (the destination). Money moves from the
// source to the destination account.
func (tr transactorAdapter) Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error {
	var customerID uuid.UUID
	switch destCustomerType {
	case models.UserTypAgent:
		agt, err := tr.customerFinder.FindAgentByEmail(destAccNumber)
		if errors.ErrorCode(err) == errors.ENOTFOUND {
			return errors.Error{Err: err, Message: errors.ErrUserNotFound}
		} else if err != nil {
			return err
		}

		customerID = agt.ID
	case models.UserTypMerchant:
		merch, err := tr.customerFinder.FindMerchantByEmail(destAccNumber)
		if errors.ErrorCode(err) == errors.ENOTFOUND {
			return errors.Error{Err: err, Message: errors.ErrUserNotFound}
		} else if err != nil {
			return err
		}

		customerID = merch.ID
	case models.UserTypeSubscriber:
		sub, err := tr.customerFinder.FindSubscriberByEmail(destAccNumber)
		if errors.ErrorCode(err) == errors.ENOTFOUND {
			return errors.Error{Err: err, Message: errors.ErrUserNotFound}
		} else if err != nil {
			return err
		}

		customerID = sub.ID
	}

	tx := transaction.Transaction{
		Source: source,
		Destination: models.TxnCustomer{
			UserID:   customerID,
			UserType: destCustomerType,
		},

		TxnOperation: models.TxnOpTransfer,
		Amount:       amount,
	}
	err := tr.transactor.Transact(tx)
	if err != nil {
		return err
	}

	return nil
}

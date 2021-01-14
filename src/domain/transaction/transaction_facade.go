package transaction

import (
	"simple-mpesa/src/domain/customer"
	"simple-mpesa/src/domain/value_objects"

	"github.com/gofrs/uuid"
)

// Facade is not a good name. Its function is to expose an interface to the application layer that it
// can use to perform transactions.
//
// To keep the Transaction context clean from a dependency of the agent, merchant and subscriber contexts,
// i chose to create this port separately.
type Facade interface {
	Deposit(depositor TxnCustomer, customerNumber string, customerType value_objects.UserType, amount value_objects.Shillings) error
	Transfer(source TxnCustomer, destAccNumber string, destCustomerType value_objects.UserType, amount value_objects.Shillings) error
	Withdraw(withdrawer TxnCustomer, agentNumber string, amount value_objects.Shillings) error
}

func NewFacade(finder customer.Finder, transactor Transactor) Facade {
	return &transactorFacade{
		customerFinder: finder,
		transactor:     transactor,
	}
}

type transactorFacade struct {
	customerFinder customer.Finder
	transactor     Transactor
}

// Deposit is a transaction between a customer and an agent. The customer's account is credited from the
// agent's account. Money moves from the agent's account to the customer's account.
// It is important to remember that it is the agent that does the deposit operation on behalf of the customer.
func (facade transactorFacade) Deposit(depositor TxnCustomer, customerNumber string, customerType value_objects.UserType, amount value_objects.Shillings) error {
	customerID, err := facade.customerFinder.FindIDByEmail(customerNumber, customerType)
	if err != nil {
		return err
	}

	tx := Transaction{
		Source: depositor,
		Destination: TxnCustomer{
			UserID:   customerID,
			UserType: customerType,
		},

		TxnOperation: value_objects.TxnOpDeposit,
		Amount:       amount,
	}
	err = facade.transactor.Transact(tx)
	if err != nil {
		return err
	}

	return nil
}

// Withdraw is a transaction between a customer and an agent. The customer's account is debited and the
// agent's account credited. Money moves from the customer's account to the agent's account.
func (facade transactorFacade) Withdraw(withdrawer TxnCustomer, agentNumber string, amount value_objects.Shillings) error {
	agt, err := facade.customerFinder.FindAgentByEmail(agentNumber)
	if err != nil {
		return err
	}

	tx := Transaction{
		Source: withdrawer,
		Destination: TxnCustomer{
			UserID:   agt.ID,
			UserType: value_objects.UserTypAgent,
		},

		TxnOperation: value_objects.TxnOpWithdraw,
		Amount:       amount,
	}
	err = facade.transactor.Transact(tx)
	if err != nil {
		return err
	}

	return nil
}

// Transfer is a transaction describing a general movement of funds from a customer to another customer. One customer's
// account is debited (the source) and the other customer's account credited (the destination). Money moves from the
// source to the destination account.
func (facade transactorFacade) Transfer(source TxnCustomer, destAccNumber string, destCustomerType value_objects.UserType, amount value_objects.Shillings) error {
	var customerID uuid.UUID
	switch destCustomerType {
	case value_objects.UserTypAgent:
		agt, err := facade.customerFinder.FindAgentByEmail(destAccNumber)
		if err != nil {
			return err
		}

		customerID = agt.ID
	case value_objects.UserTypMerchant:
		merch, err := facade.customerFinder.FindMerchantByEmail(destAccNumber)
		if err != nil {
			return err
		}

		customerID = merch.ID
	case value_objects.UserTypSubscriber:
		sub, err := facade.customerFinder.FindSubscriberByEmail(destAccNumber)
		if err != nil {
			return err
		}

		customerID = sub.ID
	}

	tx := Transaction{
		Source: source,
		Destination: TxnCustomer{
			UserID:   customerID,
			UserType: destCustomerType,
		},

		TxnOperation: value_objects.TxnOpTransfer,
		Amount:       amount,
	}
	err := facade.transactor.Transact(tx)
	if err != nil {
		return err
	}

	return nil
}

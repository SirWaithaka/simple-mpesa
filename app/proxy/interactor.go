package proxy

import (
	"simple-mpesa/app/agent"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/models"
	"simple-mpesa/app/subscriber"

	"github.com/gofrs/uuid"
)

type Interactor interface {
	Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error
	Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error
	Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error
}

func NewInteractor(agentRepo agent.Repository, merchRepo merchant.Repository, subRepo subscriber.Repository) Interactor {
	return &interactor{
		agentRepo: agentRepo,
		merchRepo: merchRepo,
		subRepo:   subRepo,
	}
}

type interactor struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository
}

func (i interactor) postTransaction(transaction models.TxnOperation) {

}

// Deposit is a transaction between a customer and an agent. The customer's account is credited from the
// agent's account. Money moves from the agent's account to the customer's account.
func (i interactor) Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := i.agentRepo.FindByEmail(agentNumber)
	if err != nil {
		return err
	}

	operation := models.TxnOperation{
		Source: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},
		Destination: depositor,

		TxnType: models.TxTypeDeposit,
		Amount:  amount,
	}

	i.postTransaction(operation)
	return nil
}

// Withdraw is a transaction between a customer and an agent. The customer's account is debited and the
// agent's account credited. Money moves from the customer's account to the agent's account.
func (i interactor) Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := i.agentRepo.FindByEmail(agentNumber)
	if err != nil {
		return err
	}

	operation := models.TxnOperation{
		Source: withdrawer,
		Destination: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},

		TxnType: models.TxTypeWithdrawal,
		Amount:  amount,
	}

	i.postTransaction(operation)
	return nil
}

// Transfer is a transaction describing a general movement of funds from a customer to another customer. One customer's
// account is debited (the source) and the other customer's account credited (the destination). Money moves from the
// source to the destination account.
func (i interactor) Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error {
	var customerID uuid.UUID
	switch destCustomerType {
	case models.UserTypAgent:
		agt, err := i.agentRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = agt.ID
	case models.UserTypMerchant:
		merch, err := i.merchRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = merch.ID
	case models.UserTypeSubscriber:
		sub, err := i.subRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = sub.ID
	}

	operation := models.TxnOperation{
		Source: source,
		Destination: models.TxnCustomer{
			UserID:   customerID,
			UserType: destCustomerType,
		},

		TxnType: models.TxTypeTransfer,
		Amount:  amount,
	}

	i.postTransaction(operation)
	return nil
}

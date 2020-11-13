package ports

import (
	"simple-mpesa/app/agent"
	"simple-mpesa/app/data"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/models"
	"simple-mpesa/app/subscriber"

	"github.com/gofrs/uuid"
)

type Transactor interface {
	Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error
	Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error
	Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error
}

func NewInteractor(agentRepo agent.Repository, merchRepo merchant.Repository, subRepo subscriber.Repository, txEventsChan data.ChanNewTxnEvents) Transactor {
	return &transactor{
		agentRepo: agentRepo,
		merchRepo: merchRepo,
		subRepo:   subRepo,
		txnEventsChannel: txEventsChan,
	}
}

type transactor struct {
	agentRepo agent.Repository
	merchRepo merchant.Repository
	subRepo   subscriber.Repository

	txnEventsChannel data.ChanNewTxnEvents
}

// posts a new transaction event to channel
func (tr transactor) postTransaction(txEvent models.TxnEvent) {
	go func() { tr.txnEventsChannel.Writer <- txEvent }()
}

// Deposit is a transaction between a customer and an agent. The customer's account is credited from the
// agent's account. Money moves from the agent's account to the customer's account.
func (tr transactor) Deposit(depositor models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := tr.agentRepo.FindByEmail(agentNumber)
	if err != nil {
		return err
	}

	event := models.TxnEvent{
		Source: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},
		Destination: depositor,

		TxnOperation: models.TxnOpDeposit,
		Amount:       amount,
	}

	tr.postTransaction(event)
	return nil
}

// Withdraw is a transaction between a customer and an agent. The customer's account is debited and the
// agent's account credited. Money moves from the customer's account to the agent's account.
func (tr transactor) Withdraw(withdrawer models.TxnCustomer, agentNumber string, amount models.Shillings) error {
	agt, err := tr.agentRepo.FindByEmail(agentNumber)
	if err != nil {
		return err
	}

	event := models.TxnEvent{
		Source: withdrawer,
		Destination: models.TxnCustomer{
			UserID:   agt.ID,
			UserType: models.UserTypAgent,
		},

		TxnOperation: models.TxnOpWithdrawal,
		Amount:       amount,
	}

	tr.postTransaction(event)
	return nil
}

// Transfer is a transaction describing a general movement of funds from a customer to another customer. One customer's
// account is debited (the source) and the other customer's account credited (the destination). Money moves from the
// source to the destination account.
func (tr transactor) Transfer(source models.TxnCustomer, destAccNumber string, destCustomerType models.UserType, amount models.Shillings) error {
	var customerID uuid.UUID
	switch destCustomerType {
	case models.UserTypAgent:
		agt, err := tr.agentRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = agt.ID
	case models.UserTypMerchant:
		merch, err := tr.merchRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = merch.ID
	case models.UserTypeSubscriber:
		sub, err := tr.subRepo.FindByEmail(destAccNumber)
		if err != nil {
			return err
		}
		customerID = sub.ID
	}

	event := models.TxnEvent{
		Source: source,
		Destination: models.TxnCustomer{
			UserID:   customerID,
			UserType: destCustomerType,
		},

		TxnOperation: models.TxnOpTransfer,
		Amount:       amount,
	}

	tr.postTransaction(event)
	return nil
}

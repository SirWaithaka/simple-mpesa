package transaction

import (
	"log"

	"simple-mpesa/app/account"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
)

type Transactor interface {
	Transact(Transaction) error
}

func NewTransactor(accountant account.Accountant) Transactor {
	return &transactor{accountant}
}

type transactor struct {
	accountant account.Accountant
}

// in mobile money a deposit will happen from the account of an agent to the other customer. The source is the agent's
// account and destination is the account of the other customer.
func (tr transactor) deposit(source, destination models.TxnCustomer, amount models.Shillings) error {
	if amount < minimumDepositAmount {
		e := errors.ErrAmountBelowMinimum(minimumDepositAmount, errors.DepositAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	// the source should always be an agent
	if source.UserType != models.UserTypAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.DepositOnlyAtAgent}
	}

	// a merchant is not allowed to deposit
	if destination.UserType == models.UserTypMerchant {
		return errors.Error{Code: errors.EINVALID, Message: errors.CustomerCantDeposit}
	}

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amount)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

// in mobile money a withdrawal will happen from the account of the customer withdrawing to the agent. The source is the
// customer's account and the destination is the account of the agent
func (tr transactor) withdraw(source, destination models.TxnCustomer, amount models.Shillings) error {
	if amount < minimumWithdrawalAmount {
		e := errors.ErrAmountBelowMinimum(minimumWithdrawalAmount, errors.WithdrawAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	// the destination should always be an agent
	if destination.UserType != models.UserTypAgent {
		return errors.Error{Code: errors.EINVALID, Message: errors.WithdrawalOnlyAtAgent}
	}

	// we can implement a double withdrawal check here. That will prevent a user from
	// withdrawing same amount twice within a stipulated time interval because of system lag.

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amount)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

func (tr transactor) transfer(source, destination models.TxnCustomer, amount models.Shillings) error {
	if amount < minimumTransferAmount {
		e := errors.ErrAmountBelowMinimum(minimumTransferAmount, errors.TransferAmountBelowMinimum)
		return errors.Error{Err: e}
	}

	srcNewBal, err := tr.accountant.DebitAccount(source.UserID, amount)
	if err != nil {
		return err
	}

	destNewBal, err := tr.accountant.CreditAccount(destination.UserID, amount)
	if err != nil {
		return err
	}

	log.Printf("Source balance: %v || Dest balance: %v", srcNewBal, destNewBal)

	return nil
}

func (tr transactor) Transact(transaction Transaction) error {
	switch transaction.TxnOperation {
	case models.TxnOpDeposit:
		return tr.deposit(transaction.Source, transaction.Destination, transaction.Amount)
	case models.TxnOpWithdrawal:
		return tr.withdraw(transaction.Source, transaction.Destination, transaction.Amount)
	case models.TxnOpTransfer:
		return tr.transfer(transaction.Source, transaction.Destination, transaction.Amount)
	}
	return errors.Error{Code: errors.EINVALID}
}

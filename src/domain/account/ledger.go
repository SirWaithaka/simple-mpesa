package account

import (
	"time"

	"simple-mpesa/src/domain/value_objects"

	"github.com/gofrs/uuid"
)

type Ledger interface {
	Record(userID uuid.UUID, acc Account, txnOp value_objects.TxnOperation, amount value_objects.Shillings, txnType TransactionType) error
}

func NewLedger(repository StatementRepository) Ledger {
	return &ledger{repository}
}

type ledger struct {
	statementRepo StatementRepository
}

func (l ledger) Record(userID uuid.UUID, acc Account, txnOp value_objects.TxnOperation, amount value_objects.Shillings, txnType TransactionType) error {
	statement := Statement{
		Operation: txnOp,
		UserID:    userID,
		AccountID: acc.ID,
		CreatedAt: time.Now(),
	}

	if txnType == TxnTypeCredit {
		statement.CreditAmount = float64(amount)
	} else if txnType == TxnTypeDebit {
		statement.DebitAmount = float64(amount)
	}

	_, err := l.statementRepo.Add(statement)
	if err != nil {
		return err
	}

	return nil
}

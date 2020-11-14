package statement

import (
	"time"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type Ledger interface {
	Record(userID uuid.UUID, acc models.Account, txnOp models.TxnOperation, amount models.Shillings, stmtType Type) error
}

func NewLedger(repository Repository) Ledger {
	return &ledger{repository}
}

type ledger struct {
	statementRepo Repository
}

func (l ledger) Record(userID uuid.UUID, acc models.Account, txnOp models.TxnOperation, amount models.Shillings, stmtType Type) error {
	statement := Statement{
		Operation:    txnOp,
		UserID:       userID,
		AccountID:    acc.ID,
		CreatedAt:    time.Now(),
	}

	if stmtType == TypeCredit {
		statement.CreditAmount = float64(amount)
	} else if stmtType == TypeDebit {
		statement.DebitAmount = float64(amount)
	}

	_, err := l.statementRepo.Add(statement)
	if err != nil {
		return err
	}

	return nil
}
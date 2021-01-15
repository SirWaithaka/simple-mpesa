package pg

import (
	"time"

	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/repositories/schema"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
)

func NewStatementRepository(database *storage.Database) *StatementRepository {
	return &StatementRepository{db: database}
}

type StatementRepository struct {
	db *storage.Database
}

func (r StatementRepository) Add(stmt account.Statement) (account.Statement, error) {
	row := schema.Statement{
		Operation:    string(stmt.Operation),
		DebitAmount:  stmt.DebitAmount,
		CreditAmount: stmt.CreditAmount,
		UserID:       stmt.UserID,
		AccountID:    stmt.AccountID,
		CreatedAt:    stmt.CreatedAt,
	}

	result := r.db.Create(&row)
	if err := result.Error; err != nil {
		return account.Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return account.Statement{
		ID:           row.ID,
		Operation:    value_objects.TxnOperation(row.Operation),
		DebitAmount:  row.DebitAmount,
		CreditAmount: row.CreditAmount,
		UserID:       row.UserID,
		AccountID:    row.AccountID,
		CreatedAt:    row.CreatedAt,
	}, nil
}

func (r StatementRepository) GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]account.Statement, error) {
	var rows []schema.Statement

	result := r.db.Where(
		schema.Statement{UserID: userID},
	).Where(
		"created_at <= ?", from,
	).Order("created_at desc").Limit(int(limit)).Find(&rows)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	var statements []account.Statement
	for _, row := range rows {
		statements = append(statements, account.Statement{
			ID:           row.ID,
			Operation:    value_objects.TxnOperation(row.Operation),
			DebitAmount:  row.DebitAmount,
			CreditAmount: row.CreditAmount,
			UserID:       row.UserID,
			AccountID:    row.AccountID,
			CreatedAt:    row.CreatedAt,
		})
	}

	return statements, nil
}

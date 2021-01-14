package pg

import (
	"time"

	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/errors"
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
	result := r.db.Create(&stmt)
	if err := result.Error; err != nil {
		return account.Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return stmt, nil
}

func (r StatementRepository) GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]account.Statement, error) {
	var statements []account.Statement

	result := r.db.Where(
		account.Statement{UserID: userID},
	).Where(
		"created_at <= ?", from,
	).Order("created_at desc").Limit(int(limit)).Find(&statements)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return statements, nil
}

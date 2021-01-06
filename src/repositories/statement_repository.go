package repositories

import (
	"time"

	"simple-mpesa/src/errors"
	"simple-mpesa/src/statement"
	"simple-mpesa/src/storage"

	"github.com/gofrs/uuid"
)

func NewStatementRepository(database *storage.Database) *Statement {
	return &Statement{db: database}
}

type Statement struct {
	db *storage.Database
}

func (r Statement) Add(stmt statement.Statement) (statement.Statement, error) {
	result := r.db.Create(&stmt)
	if err := result.Error; err != nil {
		return statement.Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return stmt, nil
}

func (r Statement) GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]statement.Statement, error) {
	var statements []statement.Statement

	result := r.db.Where(
		statement.Statement{UserID: userID},
	).Where(
		"created_at <= ?", from,
	).Order("created_at desc").Limit(int(limit)).Find(&statements)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return statements, nil
}

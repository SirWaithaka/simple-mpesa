package statement

import (
	"time"

	"simple-mpesa/app/errors"
	"simple-mpesa/app/storage"

	"github.com/gofrs/uuid"
)

type Repository interface {
	Add(Statement) (Statement, error)
	GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error)
}

func NewRepository(database *storage.Database) Repository {
	return &repository{db: database}
}

type repository struct {
	db *storage.Database
}

func (r repository) Add(stmt Statement) (Statement, error) {
	result := r.db.Create(&stmt)
	if err := result.Error; err != nil {
		return Statement{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return stmt, nil
}

func (r repository) GetStatements(userID uuid.UUID, from time.Time, limit uint) ([]Statement, error) {
	var statements []Statement

	result := r.db.Where(
		Statement{UserID: userID},
	).Where(
		"created_at <= ?", from,
	).Order("created_at desc").Limit(int(limit)).Find(&statements)

	if err := result.Error; err != nil {
		return nil, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return statements, nil
}

package statement

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/storage"
)

type Repository interface {
	Add(Statement) (Statement, error)
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

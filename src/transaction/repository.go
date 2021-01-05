package transaction

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"
)

type Repository interface {
	Add(models.Transaction) (models.Transaction, error)
}

type repository struct {
	database *storage.Database
}

func NewRepository(db *storage.Database) Repository {
	return &repository{db}
}

func (r repository) Add(tx models.Transaction) (models.Transaction, error) {
	result := r.database.Create(&tx)
	if err := result.Error; err != nil {
		return models.Transaction{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tx, nil
}

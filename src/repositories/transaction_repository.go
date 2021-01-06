package repositories

import (
	"simple-mpesa/src/errors"
	"simple-mpesa/src/models"
	"simple-mpesa/src/storage"
)

type Transaction struct {
	database *storage.Database
}

func NewTransactionRepository(db *storage.Database) *Transaction {
	return &Transaction{db}
}

func (r Transaction) Add(tx models.Transaction) (models.Transaction, error) {
	result := r.database.Create(&tx)
	if err := result.Error; err != nil {
		return models.Transaction{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tx, nil
}

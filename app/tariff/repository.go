package tariff

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(Tariff) (Tariff, error)
	Get(operation models.TxnOperation, src models.UserType, dest models.UserType) (Tariff, error)
}

func NewRepository(db *storage.Database) Repository {
	return &repository{db}
}

type repository struct {
	db *storage.Database
}

func (r repository) Add(tariff Tariff) (Tariff, error) {
	result := r.db.Create(&tariff)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return Tariff{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrTariffExists}
		}
		return Tariff{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tariff, nil
}

func (r repository) Get(operation models.TxnOperation, src models.UserType, dest models.UserType) (Tariff, error) {
	row := Tariff{Transaction: operation, SourceUserType: src, DestinationUserType: dest}

	var tariff Tariff
	result := r.db.Where(row).First(&tariff)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Tariff{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return Tariff{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tariff, nil
}

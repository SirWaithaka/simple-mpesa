package pg

import (
	"simple-mpesa/src/domain/tariff"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/storage"
	"simple-mpesa/src/value_objects"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

func NewTariffRepository(db *storage.Database) *Tariff {
	return &Tariff{db}
}

type Tariff struct {
	db *storage.Database
}

func (r Tariff) Add(charge tariff.Charge) (tariff.Charge, error) {
	result := r.db.Create(&charge)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return tariff.Charge{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrChargeExists}
		}
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r Tariff) FetchAll() ([]tariff.Charge, error) {
	var charges []tariff.Charge
	result := r.db.Find(&charges)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return charges, nil
}

func (r Tariff) FindByID(id uuid.UUID) (tariff.Charge, error) {
	var charge tariff.Charge
	result := r.db.Where(tariff.Charge{ID: id}).First(&charge)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tariff.Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r Tariff) Get(operation value_objects.TxnOperation, src value_objects.UserType, dest value_objects.UserType) (tariff.Charge, error) {
	row := tariff.Charge{Transaction: operation, SourceUserType: src, DestinationUserType: dest}

	var charge tariff.Charge
	result := r.db.Where(row).First(&charge)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tariff.Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r Tariff) Update(charge tariff.Charge) error {
	var ch tariff.Charge
	result := r.db.Model(&ch).Where(tariff.Charge{ID: charge.ID}).Updates(charge)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}


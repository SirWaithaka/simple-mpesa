package tariff

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/storage"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type Repository interface {
	Add(Charge) (Charge, error)
	FetchAll() ([]Charge, error)
	FindByID(uuid.UUID) (Charge, error)
	Get(operation models.TxnOperation, src models.UserType, dest models.UserType) (Charge, error)
	Update(Charge) error
}

func NewRepository(db *storage.Database) Repository {
	return &repository{db}
}

type repository struct {
	db *storage.Database
}

func (r repository) Add(charge Charge) (Charge, error) {
	result := r.db.Create(&charge)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return Charge{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrChargeExists}
		}
		return Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r repository) FetchAll() ([]Charge, error) {
	var charges []Charge
	result := r.db.Find(&charges)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	return charges, nil
}

func (r repository) FindByID(id uuid.UUID) (Charge, error) {
	var charge Charge
	result := r.db.Where(Charge{ID: id}).First(&charge)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r repository) Get(operation models.TxnOperation, src models.UserType, dest models.UserType) (Charge, error) {
	row := Charge{Transaction: operation, SourceUserType: src, DestinationUserType: dest}

	var charge Charge
	result := r.db.Where(row).First(&charge)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return charge, nil
}

func (r repository) Update(charge Charge) error {
	var ch Charge
	result := r.db.Model(&ch).Where(Charge{ID: charge.ID}).Updates(charge)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

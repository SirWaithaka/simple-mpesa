package pg

import (
	"simple-mpesa/src/domain/tariff"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/repositories/schema"
	"simple-mpesa/src/storage"

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
	row := schema.Charge{
		Transaction:         string(charge.Transaction),
		SourceUserType:      string(charge.SourceUserType),
		DestinationUserType: string(charge.DestinationUserType),
		Fee:                 uint(charge.Fee),
	}

	result := r.db.Create(&row)
	if err := result.Error; err != nil {
		// we check if the error is a postgres unique constraint violation
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == "23505" {
			return tariff.Charge{}, errors.Error{Code: errors.ECONFLICT, Message: errors.ErrChargeExists}
		}
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tariff.Charge{
		ID:                  row.ID,
		Transaction:         value_objects.TxnOperation(row.Transaction),
		SourceUserType:      value_objects.UserType(row.SourceUserType),
		DestinationUserType: value_objects.UserType(row.DestinationUserType),
		Fee:                 value_objects.Cents(row.Fee),
	}, nil
}

func (r Tariff) FetchAll() ([]tariff.Charge, error) {
	var rows []schema.Charge
	result := r.db.Find(&rows)
	if err := result.Error; err != nil {
		return nil, errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}

	var charges []tariff.Charge
	for _, row := range rows {
		charges = append(charges, tariff.Charge{
			ID:                  row.ID,
			Transaction:         value_objects.TxnOperation(row.Transaction),
			SourceUserType:      value_objects.UserType(row.SourceUserType),
			DestinationUserType: value_objects.UserType(row.DestinationUserType),
			Fee:                 value_objects.Cents(row.Fee),
		})
	}

	return charges, nil
}

func (r Tariff) FindByID(id uuid.UUID) (tariff.Charge, error) {
	var row schema.Charge
	result := r.db.Where(schema.Charge{ID: id}).First(&row)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tariff.Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tariff.Charge{
		ID:                  row.ID,
		Transaction:         value_objects.TxnOperation(row.Transaction),
		SourceUserType:      value_objects.UserType(row.SourceUserType),
		DestinationUserType: value_objects.UserType(row.DestinationUserType),
		Fee:                 value_objects.Cents(row.Fee),
	}, nil
}

func (r Tariff) Get(operation value_objects.TxnOperation, src value_objects.UserType, dest value_objects.UserType) (tariff.Charge, error) {
	search := schema.Charge{Transaction: string(operation), SourceUserType: string(src), DestinationUserType: string(dest)}

	var row schema.Charge
	result := r.db.Where(search).First(&row)
	// check if no record found.
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return tariff.Charge{}, errors.Error{Code: errors.ENOTFOUND}
	}
	if err := result.Error; err != nil {
		return tariff.Charge{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	return tariff.Charge{
		ID:                  row.ID,
		Transaction:         value_objects.TxnOperation(row.Transaction),
		SourceUserType:      value_objects.UserType(row.SourceUserType),
		DestinationUserType: value_objects.UserType(row.DestinationUserType),
		Fee:                 value_objects.Cents(row.Fee),
	}, nil
}

func (r Tariff) Update(charge tariff.Charge) error {
	row := schema.Charge{
		Transaction:         string(charge.Transaction),
		SourceUserType:      string(charge.SourceUserType),
		DestinationUserType: string(charge.DestinationUserType),
		Fee:                 uint(charge.Fee),
	}

	result := r.db.Model(&schema.Charge{}).Where(schema.Charge{ID: charge.ID}).Updates(row)
	if err := result.Error; err != nil {
		return errors.Error{Err: result.Error, Code: errors.EINTERNAL}
	}
	return nil
}

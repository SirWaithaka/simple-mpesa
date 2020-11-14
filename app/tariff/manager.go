package tariff

import (
	"strings"

	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
)

type Manager interface {
	AddCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType, charge models.Cents) error
	GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Cents, error)
}

func NewManager(repository Repository) Manager {
	return &manager{repository}
}

type manager struct {
	repository Repository
}

func (mg manager) AddCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType, charge models.Cents) error {
	// just in case input is not in upper case
	op := models.TxnOperation(strings.ToUpper(string(txOperation)))

	if ok := models.IsValidTxnOperation(op); !ok {
		return errors.Error{Code: errors.EINVALID, Message: errors.ErrInvalidOperation}
	}

	_, err := mg.repository.Add(Tariff{
		Transaction:         op,
		SourceUserType:      source,
		DestinationUserType: dest,
		Charge:              charge,
	})
	if err != nil {
		return err
	}

	return nil
}

func (mg manager) GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Cents, error) {
	tariff, err := mg.repository.Get(operation, src, dest)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Cents(0), errors.Error{Err: err, Message: errors.ErrTariffNotSet}
	} else if err != nil {
		return models.Cents(0), err
	}

	return tariff.Charge, nil
}

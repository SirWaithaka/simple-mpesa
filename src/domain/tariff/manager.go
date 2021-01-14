package tariff

import (
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"

	"github.com/gofrs/uuid"
)

type Manager interface {
	GetCharge(operation value_objects.TxnOperation, src value_objects.UserType, dest value_objects.UserType) (value_objects.Cents, error)
	GetTariff() ([]Charge, error)
	UpdateCharge(chargeID uuid.UUID, fee value_objects.Cents) error
}

func NewManager(repository Repository) Manager {
	mgr := &manager{repository}

	go mgr.initTariffSetup()

	return mgr
}

type manager struct {
	repository Repository
}

func (mg manager) GetCharge(operation value_objects.TxnOperation, src value_objects.UserType, dest value_objects.UserType) (value_objects.Cents, error) {
	tariff, err := mg.repository.Get(operation, src, dest)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return value_objects.Cents(0), errors.Error{Err: err, Message: errors.ErrTariffNotSet}
	} else if err != nil {
		return value_objects.Cents(0), err
	}

	return tariff.Fee, nil
}

func (mg manager) GetTariff() ([]Charge, error) {
	charges, err := mg.repository.FetchAll()
	if err != nil {
		return nil, err
	}

	if len(charges) == 0 {
		return nil, errors.Error{Code: errors.ENOTFOUND, Message: errors.ErrTariffNotSet}
	}

	return charges, nil
}

func (mg manager) UpdateCharge(chargeID uuid.UUID, fee value_objects.Cents) error {
	charge, err := mg.repository.FindByID(chargeID)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return errors.Error{Err: err, Message: errors.ErrChargeNotFound}
	} else if err != nil {
		return err
	}

	charge.Fee = fee
	err = mg.repository.Update(charge)
	if err != nil {
		return err
	}

	return nil
}

// initializes a tariff with zero amount, is used only once during initial setup of charges
func (mg manager) addCharge(txOperation value_objects.TxnOperation, source value_objects.UserType, dest value_objects.UserType) error {

	if !txOperation.IsPrimaryOperation() {
		return errors.Error{Code: errors.EINVALID, Message: errors.ErrInvalidOperation}
	}

	_, err := mg.repository.Add(Charge{
		Transaction:         txOperation,
		SourceUserType:      source,
		DestinationUserType: dest,
		Fee:                 value_objects.Cents(0),
	})
	if err != nil {
		return err
	}

	return nil
}

// is a definition of all valid withdrawals, remember withdrawals can only happen
// at an agent's desk
func (mg manager) validWithdrawTx() []ValidTransaction {
	return []ValidTransaction{
		{value_objects.UserTypSubscriber, value_objects.UserTypAgent},
		{value_objects.UserTypMerchant, value_objects.UserTypAgent},
		{value_objects.UserTypAgent, value_objects.UserTypAgent},
	}
}

// is a definition of all valid transfers, remember only agents are allowed
// to transfer to other agents
func (mg manager) validTransferTx() []ValidTransaction {
	return []ValidTransaction{
		{value_objects.UserTypAgent, value_objects.UserTypAgent},           // transfer between an agent to an agent
		{value_objects.UserTypSubscriber, value_objects.UserTypSubscriber}, // transfer between a subscriber to subscriber
		{value_objects.UserTypMerchant, value_objects.UserTypSubscriber},   // transfer between a merchant to subscriber
		{value_objects.UserTypSubscriber, value_objects.UserTypMerchant},   // transfer between a subscriber to merchant -> PAYMENT
		{value_objects.UserTypAgent, value_objects.UserTypMerchant},        // transfer between an agent to merchant -> PAYMENT
	}
}

func (mg manager) initTariffSetup() error {
	// add valid withdraw transactions between customers
	for _, validTx := range mg.validWithdrawTx() {
		err := mg.addCharge(value_objects.TxnOpWithdraw, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	// add valid transfer transactions between customers
	for _, validTx := range mg.validTransferTx() {
		err := mg.addCharge(value_objects.TxnOpTransfer, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	return nil
}

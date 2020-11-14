package tariff

import (
	"strings"

	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
)

type Manager interface {
	GetCharge(operation models.TxnOperation, src models.UserType, dest models.UserType) (models.Cents, error)
	UpdateCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType, charge models.Cents) error
}

func NewManager(repository Repository) Manager {
	mgr := &manager{repository}

	go mgr.initTariffSetup()

	return mgr
}

type manager struct {
	repository Repository
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

func (mg manager) UpdateCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType, charge models.Cents) error {
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

// initializes a tariff with zero amount, is used only once during initial setup of charges
func (mg manager) addCharge(txOperation models.TxnOperation, source models.UserType, dest models.UserType) error {

	if ok := models.IsValidTxnOperation(txOperation); !ok {
		return errors.Error{Code: errors.EINVALID, Message: errors.ErrInvalidOperation}
	}

	_, err := mg.repository.Add(Tariff{
		Transaction:         txOperation,
		SourceUserType:      source,
		DestinationUserType: dest,
		Charge:              models.Cents(0),
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
		{models.UserTypSubscriber, models.UserTypAgent},
		{models.UserTypMerchant, models.UserTypAgent},
		{models.UserTypAgent, models.UserTypAgent},
	}
}

// is a definition of all valid transfers, remember only agents are allowed
// to transfer to other agents
func (mg manager) validTransferTx() []ValidTransaction {
	return []ValidTransaction{
		{models.UserTypAgent, models.UserTypAgent},           // transfer between an agent to an agent
		{models.UserTypSubscriber, models.UserTypSubscriber}, // transfer between a subscriber to subscriber
		{models.UserTypMerchant, models.UserTypSubscriber},   // transfer between a merchant to subscriber
		{models.UserTypSubscriber, models.UserTypMerchant},   // transfer between a subscriber to merchant -> PAYMENT
		{models.UserTypAgent, models.UserTypMerchant},        // transfer between an agent to merchant -> PAYMENT
	}
}

func (mg manager) initTariffSetup() error {
	// add valid withdraw transactions between customers
	for _, validTx := range mg.validWithdrawTx() {
		err := mg.addCharge(models.TxnOpWithdraw, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	// add valid transfer transactions between customers
	for _, validTx := range mg.validTransferTx() {
		err := mg.addCharge(models.TxnOpTransfer, validTx[0], validTx[1])
		if err != nil {
			return err
		}
	}

	return nil
}

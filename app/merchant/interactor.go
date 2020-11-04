package merchant

import (
	"simple-mpesa/app"
	"simple-mpesa/app/data"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/helpers"
	"simple-mpesa/app/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Merchant, error)
	Register(RegistrationParams) (models.Merchant, error)
}

func NewInteractor(config app.Config, merchRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       merchRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           app.Config
	repository       Repository
}

// AuthenticateByEmail verifies a merchant by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Merchant, error) {
	// search for merchant by email.
	merchant, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Merchant{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Merchant{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(merchant.Password, password); err != nil {
		return models.Merchant{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return merchant, nil
}

// Register takes in a merchant registration parameters and creates a new merchant
// then adds the merchant to db.
func (ui interactor) Register(params RegistrationParams) (models.Merchant, error) {
	merchant := models.Merchant{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash merchant password before adding to db.
	passwordHash, err := helpers.HashPassword(merchant.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Merchant{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	merchant.Password = passwordHash
	merch, err := ui.repository.Add(merchant)
	if err != nil {
		return models.Merchant{}, err
	}

	// tell channel listeners that a new merchant has been created.
	ui.postNewMerchantToChannel(&merch)
	return merch, nil
}

// take the newly created merchant and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewMerchantToChannel(merchant *models.Merchant) {
	newMerchant := parseToNewMerchant(*merchant)
	go func() { ui.customersChannel.Writer <- newMerchant }()
}

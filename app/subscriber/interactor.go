package subscriber

import (
	"simple-mpesa/app"
	"simple-mpesa/app/data"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/helpers"
	"simple-mpesa/app/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Subscriber, error)
	Register(RegistrationParams) (models.Subscriber, error)
}

func NewInteractor(config app.Config, subsRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       subsRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           app.Config
	repository       Repository
}

// AuthenticateByEmail verifies a subscriber by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Subscriber, error) {
	// search for subscriber by email.
	subscriber, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Subscriber{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(subscriber.Password, password); err != nil {
		return models.Subscriber{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return subscriber, nil
}

// Register takes in a subscriber registration parameters and creates a new subscriber
// then adds the subscriber to db.
func (ui interactor) Register(params RegistrationParams) (models.Subscriber, error) {
	subscriber := models.Subscriber{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash subscriber password before adding to db.
	passwordHash, err := helpers.HashPassword(subscriber.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Subscriber{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	subscriber.Password = passwordHash
	sub, err := ui.repository.Add(subscriber)
	if err != nil {
		return models.Subscriber{}, err
	}

	// tell channel listeners that a new subscriber has been created.
	ui.postNewSubscriberToChannel(&sub)
	return sub, nil
}

// take the newly created subscriber and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewSubscriberToChannel(subscriber *models.Subscriber) {
	newSubscriber := parseToNewSubscriber(*subscriber)
	go func() { ui.customersChannel.Writer <- newSubscriber }()
}

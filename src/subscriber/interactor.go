package subscriber

import (
	"simple-mpesa/src"
	"simple-mpesa/src/data"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/helpers"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (Subscriber, error)
	Register(RegistrationParams) (Subscriber, error)
}

func NewInteractor(config src.Config, subsRepo Repository, custChan data.ChanNewCustomers) Interactor {
	return &interactor{
		config:           config,
		repository:       subsRepo,
		customersChannel: custChan,
	}
}

type interactor struct {
	customersChannel data.ChanNewCustomers
	config           src.Config
	repository       Repository
}

// AuthenticateByEmail verifies a subscriber by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (Subscriber, error) {
	// search for subscriber by email.
	subscriber, err := ui.repository.FindByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return Subscriber{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return Subscriber{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(subscriber.Password, password); err != nil {
		return Subscriber{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return subscriber, nil
}

// Register takes in a subscriber registration parameters and creates a new subscriber
// then adds the subscriber to db.
func (ui interactor) Register(params RegistrationParams) (Subscriber, error) {
	subscriber := Subscriber{
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
		return Subscriber{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	subscriber.Password = passwordHash
	sub, err := ui.repository.Add(subscriber)
	if err != nil {
		return Subscriber{}, err
	}

	// tell channel listeners that a new subscriber has been created.
	ui.postNewSubscriberToChannel(&sub)
	return sub, nil
}

// take the newly created subscriber and post them to channel
// that listens for newly created customers and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewSubscriberToChannel(subscriber *Subscriber) {
	newSubscriber := parseToNewSubscriber(*subscriber)
	go func() { ui.customersChannel.Writer <- newSubscriber }()
}

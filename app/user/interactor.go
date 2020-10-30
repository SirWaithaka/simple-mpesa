package user

import (
	"simple-wallet/app"
	"simple-wallet/app/data"
	"simple-wallet/app/errors"
	"simple-wallet/app/helpers"
	"simple-wallet/app/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.User, error)
	AuthenticateByPhoneNumber(email, password string) (models.User, error)
	Register(RegistrationParams) (models.User, error)
}

func NewInteractor(config app.Config, userRepo Repository, usersChan data.ChanNewUsers) Interactor {
	return &interactor{
		config:      config,
		repository:  userRepo,
		userChannel: usersChan,
	}
}

type interactor struct {
	userChannel data.ChanNewUsers
	config      app.Config
	repository  Repository
}

// AuthenticateByEmail verifies a user by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.User, error) {
	// search for user by email.
	user, err := ui.repository.GetByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.User{}, errors.Error{Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.User{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(user.Password, password); err != nil {
		return models.User{}, errors.Unauthorized{Message: err.Error()}
	}

	return user, nil
}

// AuthenticateByPhoneNumber verifies a user by the provided unique phone number
func (ui interactor) AuthenticateByPhoneNumber(phoneNo, password string) (models.User, error) {
	// search for user by phone number.
	user, err := ui.repository.GetByPhoneNumber(phoneNo)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.User{}, errors.Error{Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.User{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(user.Password, password); err != nil {
		return models.User{}, errors.Unauthorized{Message: err.Error()}
	}

	return user, nil
}

// Register takes in a user object and adds the user to db.
func (ui interactor) Register(params RegistrationParams) (models.User, error) {
	user := models.User{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
		PassportNo:  params.PassportNo,
	}

	// hash user password before adding to db.
	passwordHash, err := helpers.HashPassword(user.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.User{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	user.Password = passwordHash
	u, err := ui.repository.Add(user)
	if err != nil {
		return models.User{}, err
	}

	// tell channel listeners that a new user has been created.
	ui.postNewUserToChannel(&u)
	return u, nil
}

// take the newly created user and post them to channel
// that listens for newly created user and acts upon them
// like creating an account for them automatically.
func (ui interactor) postNewUserToChannel(user *models.User) {
	newUser := parseToNewUser(*user)
	go func() { ui.userChannel.Writer <- newUser }()
}

package admin

import (
	"simple-mpesa/app"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/helpers"
	"simple-mpesa/app/models"
)

type Interactor interface {
	AuthenticateByEmail(email, password string) (models.Admin, error)
	Register(RegistrationParams) (models.Admin, error)
}

func NewInteractor(config app.Config, adminsRepo Repository) Interactor {
	return &interactor{
		config:     config,
		repository: adminsRepo,
	}
}

type interactor struct {
	config     app.Config
	repository Repository
}

// AuthenticateByEmail verifies a admin by the provided unique email address
func (ui interactor) AuthenticateByEmail(email, password string) (models.Admin, error) {
	// search for admin by email.
	admin, err := ui.repository.GetByEmail(email)
	if errors.ErrorCode(err) == errors.ENOTFOUND {
		return models.Admin{}, errors.Error{Err: err, Message: errors.ErrUserNotFound}
	} else if err != nil {
		return models.Admin{}, err
	}

	// validate password
	if err := helpers.ComparePasswordToHash(admin.Password, password); err != nil {
		return models.Admin{}, errors.Unauthorized{Message: errors.ErrorMessage(err)}
	}

	return admin, nil
}


// Register takes in a admin object and adds the admin to db.
func (ui interactor) Register(params RegistrationParams) (models.Admin, error) {
	admin := models.Admin{
		FirstName:   params.FirstName,
		LastName:    params.LastName,
		Email:       params.Email,
		Password:    params.Password,
	}

	// hash admin password before adding to db.
	passwordHash, err := helpers.HashPassword(admin.Password)
	if err != nil { // if we get an error, it means our hashing func dint work
		return models.Admin{}, errors.Error{Err: err, Code: errors.EINTERNAL}
	}

	// change password to hashed string
	admin.Password = passwordHash
	adm, err := ui.repository.Add(admin)
	if err != nil {
		return models.Admin{}, err
	}

	return adm, nil
}

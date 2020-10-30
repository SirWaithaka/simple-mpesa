package subscriber

import (
	"simple-wallet/app/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RegistrationParams are properties required during registration of a new subscriber
type RegistrationParams struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	PassportNo  string `json:"passportNumber"`
	Password    string `json:"password"`
}

func (req RegistrationParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired))),
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.FirstName, validation.Required.Error(string(errors.ErrorFirstNameRequired))),
		validation.Field(&req.LastName, validation.Required.Error(string(errors.ErrorLastNameRequired))),
	)

	e := errors.ParseValidationErrorMap(err)
	if len(e) > 0 {
		// we will return only the first error
		return e[0]
	}
	return nil
}

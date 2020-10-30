package user

import (
	"simple-wallet/app/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// LoginParams are properties required during login of a user
// email and phoneNumber can be interchanged, either is required
type LoginParams struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
}

func (req LoginParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		// when validating email, if no phone number is provided, make sure email is valid and not empty
		validation.Field(&req.Email,
			validation.When(req.PhoneNumber == "", validation.Required.Error(string(errors.ErrorEmailRequired)), is.Email),
		),
		// when validating phone number, if no email is provided, make sure phone number is not empty
		validation.Field(&req.PhoneNumber,
			validation.When(req.Email == "", validation.Required.Error(string(errors.ErrorPhoneNumberRequired))),
		),
	)

	e := errors.ParseValidationErrorMap(err)
	if len(e) > 0 {
		// we will return only the first error
		return e[0]
	}
	return nil
}

// RegistrationParams are properties required during registration of a new user
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

package agent

import (
	"simple-mpesa/src/errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)


// LoginParams are properties required during login of an agent
type LoginParams struct {
	Email    string `json:"email" schema:"email" form:"email"`
	Password string `json:"password" schema:"password" form:"password"`
}

func (req LoginParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.EmailFormat),
	)

	return errors.ParseValidationErrorMap(err)
}

// RegistrationParams are properties required during registration of a new agent
type RegistrationParams struct {
	FirstName   string `json:"firstName" schema:"firstName" form:"firstName"`
	LastName    string `json:"lastName" schema:"lastName" form:"lastName"`
	Email       string `json:"email" schema:"email" form:"email"`
	PhoneNumber string `json:"phoneNumber" schema:"phoneNumber" form:"phoneNumber"`
	PassportNo  string `json:"passportNumber" schema:"passportNumber" form:"passportNumber"`
	Password    string `json:"password" schema:"password" form:"password"`
}

func (req RegistrationParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.EmailFormat),
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.FirstName, validation.Required.Error(string(errors.ErrorFirstNameRequired))),
		validation.Field(&req.LastName, validation.Required.Error(string(errors.ErrorLastNameRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type MakeSuperAgentParams struct {
	Email string `json:"email" schema:"email" form:"email"`
}

func (req MakeSuperAgentParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.EmailFormat),
	)

	return errors.ParseValidationErrorMap(err)
}

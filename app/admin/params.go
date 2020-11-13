package admin

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// LoginParams are properties required during login of an admin
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired)), is.Email),
	)

	return errors.ParseValidationErrorMap(err)
}

// RegistrationParams are properties required during registration of a new admin
type RegistrationParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (req RegistrationParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Email, validation.Required.Error(string(errors.ErrorEmailRequired))),
		validation.Field(&req.Password, validation.Required.Error(string(errors.ErrorPasswordRequired))),
		validation.Field(&req.FirstName, validation.Required.Error(string(errors.ErrorFirstNameRequired))),
		validation.Field(&req.LastName, validation.Required.Error(string(errors.ErrorLastNameRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type AssignFloatParams struct {
	AgentAccountNumber string           `json:"accountNo" schema:"accountNo" form:"accountNo"`
	Amount             models.Shillings `json:"amount" schema:"amount" form:"amount"`
}

func (req AssignFloatParams) Validate() error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.AgentAccountNumber, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

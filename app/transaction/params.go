package transaction

import (
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DepositParams struct {
	Amount models.Shillings `json:"amount" schema:"amount" form:"amount"`
	// In a production system, the customer number is usually a generated number for the case
	// of a merchant of agent and a mobile number for a subscriber , but we are going to use the
	// customer's email as a replacement
	CustomerNumber string          `json:"accountNo" schema:"accountNo" form:"accountNo"`
	CustomerType   models.UserType `json:"customerType" schema:"customerType" form:"customerType"`
}

func (req DepositParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.CustomerNumber, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.CustomerType, validation.Required.Error(string(errors.ErrorCustomerTypeRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type TransferParams struct {
	Amount models.Shillings `json:"amount" schema:"amount" form:"amount"`

	// In a production system, the account number is usually
	// a generated number, but we are going to use the customer's
	// email as a replacement
	DestAccountNo string          `json:"accountNo" schema:"accountNo" form:"accountNo"`
	DestUserType  models.UserType `json:"customerType" schema:"customerType" form:"customerType"`
}

func (req TransferParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.DestAccountNo, validation.Required.Error(string(errors.ErrorAccountNumberRequired))),
		validation.Field(&req.DestUserType, validation.Required.Error(string(errors.ErrorCustomerTypeRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type WithdrawParams struct {
	Amount models.Shillings `json:"amount" schema:"amount" form:"amount"`
	// In a production system, the agent number is usually
	// a generated number, but we are going to use the agent's
	// email as a replacement
	AgentNumber string `json:"agentNumber" schema:"agentNumber" form:"agentNumber"`
}

func (req WithdrawParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.AgentNumber, validation.Required.Error(string(errors.ErrorAgentNumberRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

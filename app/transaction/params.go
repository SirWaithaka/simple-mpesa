package transaction

import (
	"log"

	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type DepositParams struct {
	Amount models.Shillings `json:"amount"`
	// In a production system, the agent number is usually
	// a generated number, but we are going to use the agent's
	// email as a replacement
	AgentNumber string `json:"agentNumber"`
}

func (req DepositParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.AgentNumber, validation.Required.Error(string(errors.ErrorAgentNumberRequired))),
	)

	e := errors.ParseValidationErrorMap(err)
	log.Println(e)

	return e
}

type TransferParams struct {
	Amount models.Shillings `json:"amount"`

	// In a production system, the account number is usually
	// a generated number, but we are going to use the agent's
	// email as a replacement
	DestAccountNo string `json:"destinationAccNo"`
}

func (req TransferParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.DestAccountNo, validation.Required.Error(string(errors.ErrorDestAccNumberRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

type WithdrawParams struct {
	Amount models.Shillings `json:"amount"`
	// In a production system, the agent number is usually
	// a generated number, but we are going to use the agent's
	// email as a replacement
	AgentNumber string `json:"agentNumber"`
}

func (req WithdrawParams) Validate() error {

	err := validation.ValidateStruct(&req,
		validation.Field(&req.Amount, validation.Required.Error(string(errors.ErrorAmountRequired))),
		validation.Field(&req.AgentNumber, validation.Required.Error(string(errors.ErrorAgentNumberRequired))),
	)

	return errors.ParseValidationErrorMap(err)
}

package transaction_handlers

import (
	"net/http"

	"simple-mpesa/src/domain/auth"
	"simple-mpesa/src/domain/transaction"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/routing/responses"

	"github.com/gofiber/fiber/v2"
)

// Deposit allows user to deposit or credit their account.
func Deposit(transactor transaction.Facade) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.DepositParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		depositor := transaction.TxnCustomer{
			UserType: userDetails.UserType,
			UserID:   userDetails.UserID,
		}
		err = transactor.Deposit(depositor, p.CustomerNumber, p.CustomerType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

// Withdraw allows user to withdraw or debit their account.
func Withdraw(transactor transaction.Facade) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.WithdrawParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		withdrawer := transaction.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = transactor.Withdraw(withdrawer, p.AgentNumber, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

func Transfer(transactor transaction.Facade) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// inflate struct with body params
		var p transaction.TransferParams
		_ = ctx.BodyParser(&p)

		// validate params
		err := p.Validate()
		if err != nil {
			return err
		}

		source := transaction.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = transactor.Transfer(source, p.DestAccountNo, p.DestUserType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

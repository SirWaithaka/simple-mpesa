package transaction_handlers

import (
	"net/http"

	"simple-mpesa/app/auth"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/ports"
	"simple-mpesa/app/routing/responses"
	"simple-mpesa/app/transaction"

	"github.com/gofiber/fiber/v2"
)

// Deposit allows user to deposit or credit their account.
func Deposit(txnAdapter ports.TransactorPort) fiber.Handler {

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

		depositor := models.TxnCustomer{
			UserType: userDetails.UserType,
			UserID:   userDetails.UserID,
		}
		err = txnAdapter.Deposit(depositor, p.CustomerNumber, p.CustomerType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

// Withdraw allows user to withdraw or debit their account.
func Withdraw(txnAdapter ports.TransactorPort) fiber.Handler {

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

		withdrawer := models.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = txnAdapter.Withdraw(withdrawer, p.AgentNumber, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

func Transfer(txnAdapter ports.TransactorPort) fiber.Handler {

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

		source := models.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err = txnAdapter.Transfer(source, p.DestAccountNo, p.DestUserType, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

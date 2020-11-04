package account_handlers

import (
	"net/http"

	"simple-mpesa/app/account"
	"simple-mpesa/app/auth"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/proxy"
	"simple-mpesa/app/routing/responses"
	"simple-mpesa/app/transaction"

	"github.com/gofiber/fiber/v2"
)

// BalanceEnquiry ...
func BalanceEnquiry(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		// we check if user is admin, we return error
		if userDetails.UserType == models.UserTypAdmin {
			return errors.Error{Code: errors.EINVALID, Message: errors.UserCantHaveAccount}
		}

		balance, err := interactor.GetBalance(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.BalanceResponse(userDetails.UserID, balance))
	}
}

// Deposit allows user to deposit or credit their
// account.
func Deposit(interactor proxy.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		var p account.DepositParams
		_ = ctx.BodyParser(&p)

		depositor := models.TxnCustomer{
			UserType: userDetails.UserType,
			UserID: userDetails.UserID,
		}
		err := interactor.Deposit(depositor, p.AgentNumber, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

// Withdraw allows user to withdraw or debit their
// account.
func Withdraw(interactor proxy.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		var p account.WithdrawParams
		_ = ctx.BodyParser(&p)

		withdrawer := models.TxnCustomer{
			UserID:   userDetails.UserID,
			UserType: userDetails.UserType,
		}
		err := interactor.Withdraw(withdrawer, p.AgentNumber, p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse())
	}
}

// MiniStatement returns a small short summary of the
// most recent transactions on an account.
func MiniStatement(interactor transaction.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		transactions, err := interactor.GetStatement(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.MiniStatementResponse(userDetails.UserID, *transactions))
	}
}

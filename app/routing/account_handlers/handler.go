package account_handlers

import (
	"net/http"

	"simple-mpesa/app/account"
	"simple-mpesa/app/auth"
	"simple-mpesa/app/errors"
	"simple-mpesa/app/models"
	"simple-mpesa/app/routing/responses"
	"simple-mpesa/app/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type param struct {
	Amount uint `json:"amount"`
}

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

		balance, err := interactor.GetBalance(uuid.FromStringOrNil(userDetails.UserId))
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.BalanceResponse(userDetails.UserId, balance))
	}
}

// Deposit allows user to deposit or credit their
// account.
func Deposit(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		var p param
		_ = ctx.BodyParser(&p)

		balance, err := interactor.Deposit(uuid.FromStringOrNil(userDetails.UserId), p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse(models.TxTypeDeposit, userDetails.UserId, balance))
	}
}

// Withdraw allows user to withdraw or debit their
// account.
func Withdraw(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		var p param
		_ = ctx.BodyParser(&p)

		balance, err := interactor.Withdraw(uuid.FromStringOrNil(userDetails.UserId), p.Amount)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.TransactionResponse(models.TxTypeWithdrawal, userDetails.UserId, balance))
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

		transactions, err := interactor.GetStatement(uuid.FromStringOrNil(userDetails.UserId))
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.MiniStatementResponse(userDetails.UserId, *transactions))
	}
}

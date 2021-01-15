package account_handlers

import (
	"net/http"

	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/auth"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/errors"
	"simple-mpesa/src/routing/responses"

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
		if userDetails.UserType == value_objects.UserTypAdmin {
			return errors.Error{Code: errors.EINVALID, Message: errors.UserCantHaveAccount}
		}

		balance, err := interactor.GetBalance(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.BalanceResponse(userDetails.UserID, balance))
	}
}

// MiniStatement returns a small short summary of the
// most recent transactions on an account.
func MiniStatement(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails auth.UserAuthDetails
		if details, ok := ctx.Locals("userDetails").(auth.UserAuthDetails); !ok {
			return errors.Error{Code: errors.EINTERNAL}
		} else {
			userDetails = details
		}

		statements, err := interactor.GetStatement(userDetails.UserID)
		if err != nil {
			return err
		}

		return ctx.Status(http.StatusOK).JSON(responses.MiniStatementResponse(userDetails.UserID, statements))
	}
}

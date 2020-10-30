package account_handlers

import (
	"fmt"

	"simple-wallet/app/account"
	"simple-wallet/app/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type param struct {
	Amount uint `json:"amount"`
}

// BalanceEnquiry ...
func BalanceEnquiry(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails = ctx.Locals("userDetails").(map[string]string)
		userId := userDetails["userId"]

		balance, err := interactor.GetBalance(uuid.FromStringOrNil(userId))
		if err != nil {
			return err
		}

		return ctx.JSON(map[string]interface{}{
			"message": fmt.Sprintf("Your current balance is %v", balance),
			"balance": balance,
		})
	}
}

// Deposit allows user to deposit or credit their
// account.
func Deposit(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails = ctx.Locals("userDetails").(map[string]string)
		userId := userDetails["userId"]

		var p param
		_ = ctx.BodyParser(&p)

		balance, err := interactor.Deposit(uuid.FromStringOrNil(userId), p.Amount)
		if err != nil {
			return err
		}

		return ctx.JSON(map[string]interface{}{
			"message": fmt.Sprintf("Amount successfully deposited. New balance %v", balance),
			"balance": balance,
			"userId":  userId,
		})
	}
}

// Withdraw allows user to withdraw or debit their
// account.
func Withdraw(interactor account.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails = ctx.Locals("userDetails").(map[string]string)
		userId := userDetails["userId"]

		var p param
		_ = ctx.BodyParser(&p)

		balance, err := interactor.Withdraw(uuid.FromStringOrNil(userId), p.Amount)
		if err != nil {
			return err
		}

		return ctx.JSON(map[string]interface{}{
			"message": fmt.Sprintf("Amount successfully withdrawn. New balance %v", balance),
			"balance": balance,
			"userId":  userId,
		})
	}
}

// MiniStatement returns a small short summary of the
// most recent transactions on an account.
func MiniStatement(interactor transaction.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var userDetails = ctx.Locals("userDetails").(map[string]string)
		userId := userDetails["userId"]

		transactions, err := interactor.GetStatement(uuid.FromStringOrNil(userId))
		if err != nil {
			return err
		}

		return ctx.JSON(map[string]interface{}{
			"message":      "ministatement retrieved for the past 5 transactions",
			"userId":       userId,
			"transactions": transactions,
		})
	}
}

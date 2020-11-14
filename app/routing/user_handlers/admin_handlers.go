package user_handlers

import (
	"net/http"

	"simple-mpesa/app"
	"simple-mpesa/app/admin"
	"simple-mpesa/app/auth"
	"simple-mpesa/app/models"
	"simple-mpesa/app/routing/responses"
	"simple-mpesa/app/tariff"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

func AuthenticateAdmin(adminDomain admin.Interactor, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params admin.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		// authenticate by email.
		adm, err := adminDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating admin.
		if err != nil {
			return err
		}

		// generate an auth token string
		token, err := auth.GetTokenString(adm.ID, models.UserTypAdmin, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   adm.ID.String(),
			UserType: models.UserTypAdmin,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterAdmin(adminDomain admin.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		// register admin
		adm, err := adminDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of admin.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(adm.ID, models.UserTypAdmin))

		return nil
	}
}

func AssignFloat(adminDomain admin.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.AssignFloatParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		balance, err := adminDomain.AssignFloat(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of admin.
		_ = ctx.Status(http.StatusOK).JSON(responses.BalanceResponse(uuid.Nil, balance))

		return nil
	}
}

func AddCharge(manager tariff.Manager) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params admin.AddChargeParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return err
		}

		err = manager.AddCharge(params.TxnOperation, params.SrcUserType, params.DestUserType, params.Amount)
		if err != nil {
			return err
		}

		_ = ctx.Status(http.StatusOK).JSON(responses.SuccessResponse{
			Status:  "success",
			Message: "charge configured",
		})

		return nil
	}
}

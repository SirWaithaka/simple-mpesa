package user_handlers

import (
	"net/http"

	"simple-wallet/app"
	"simple-wallet/app/auth"
	"simple-wallet/app/merchant"
	"simple-wallet/app/models"
	"simple-wallet/app/routing/responses"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateMerchant(merchDomain merchant.Interactor, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params merchant.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// authenticate by email.
		merch, err := merchDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating merchant.
		if err != nil {
			return err
		}

		// generate an auth token string
		token, err := auth.GetTokenString(merch.ID, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   merch.ID.String(),
			UserType: models.UserTypMerchant,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterMerchant(merchDomain merchant.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params merchant.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// register merchant
		merch, err := merchDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of merchant.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(merch.ID, models.UserTypMerchant))

		return nil
	}
}

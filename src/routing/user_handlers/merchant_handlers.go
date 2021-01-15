package user_handlers

import (
	"net/http"

	"simple-mpesa/src"
	"simple-mpesa/src/domain/auth"
	"simple-mpesa/src/domain/merchant"
	"simple-mpesa/src/domain/value_objects"
	"simple-mpesa/src/routing/responses"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateMerchant(merchDomain merchant.Interactor, config src.Config) fiber.Handler {

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
		token, err := auth.GetTokenString(merch.ID, value_objects.UserTypMerchant, config.Secret)
		if err != nil {
			return err
		}

		signedUser := value_objects.SignedUser{
			UserID:   merch.ID.String(),
			UserType: value_objects.UserTypMerchant,
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
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(merch.ID, value_objects.UserTypMerchant))

		return nil
	}
}

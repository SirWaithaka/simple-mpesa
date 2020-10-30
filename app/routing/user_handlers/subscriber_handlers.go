package user_handlers

import (
	"net/http"

	"simple-wallet/app"
	"simple-wallet/app/auth"
	"simple-wallet/app/models"
	"simple-wallet/app/routing/responses"
	"simple-wallet/app/subscriber"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateSubscriber(subDomain subscriber.Interactor, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params subscriber.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// authenticate by email.
		sub, err := subDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating subscriber.
		if err != nil {
			return err
		}

		// generate an auth token string
		token, err := auth.GetTokenString(sub.ID, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   sub.ID.String(),
			UserType: models.UserTypeSubscriber,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterSubscriber(subDomain subscriber.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params subscriber.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// register subscriber
		sub, err := subDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of subscriber.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(sub.ID, models.UserTypeSubscriber))

		return nil
	}
}

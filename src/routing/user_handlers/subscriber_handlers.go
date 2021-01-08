package user_handlers

import (
	"net/http"

	"simple-mpesa/src"
	"simple-mpesa/src/auth"
	"simple-mpesa/src/routing/responses"
	"simple-mpesa/src/subscriber"
	"simple-mpesa/src/value_objects"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateSubscriber(subDomain subscriber.Interactor, config src.Config) fiber.Handler {

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
		token, err := auth.GetTokenString(sub.ID, value_objects.UserTypSubscriber, config.Secret)
		if err != nil {
			return err
		}

		signedUser := value_objects.SignedUser{
			UserID:   sub.ID.String(),
			UserType: value_objects.UserTypSubscriber,
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
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(sub.ID, value_objects.UserTypSubscriber))

		return nil
	}
}

package user_handlers

import (
	"net/http"

	"simple-wallet/app"
	"simple-wallet/app/auth"
	"simple-wallet/app/models"
	"simple-wallet/app/user"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(userDomain user.Interactor, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params user.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// user we are authenticating
		var u models.User
		switch {
		case params.Email != "":
			// if email parameter not empty authenticate by email.
			u, err = userDomain.AuthenticateByEmail(params.Email, params.Password)
		case params.PhoneNumber != "":
			// if phone number parameter not empty authenticate by phone number.
			u, err = userDomain.AuthenticateByPhoneNumber(params.PhoneNumber, params.Password)
		}

		// if there is an error authenticating user.
		if err != nil {
			return err
		}

		// generate an auth token string
		token, err := auth.GetTokenString(u.ID, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID: u.ID.String(),
			Token:  token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func Register(userDomain user.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params user.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// register user
		u, err := userDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of user.
		_ = ctx.JSON(user.RegistrationResponse(&u))

		return nil
	}
}

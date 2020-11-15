package user_handlers

import (
	"net/http"

	"simple-mpesa/app"
	"simple-mpesa/app/agent"
	"simple-mpesa/app/auth"
	"simple-mpesa/app/models"
	"simple-mpesa/app/routing/responses"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateAgent(agentDomain agent.Interactor, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		var params agent.LoginParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// authenticate by email.
		agt, err := agentDomain.AuthenticateByEmail(params.Email, params.Password)

		// if there is an error authenticating agent.
		if err != nil {
			return err
		}

		// check if agent is a super agent
		agentType := models.UserTypAgent
		if agt.IsSuperAgent() {
			agentType = models.UserTypSuperAgent
		}

		// generate an auth token string
		token, err := auth.GetTokenString(agt.ID, agentType, config.Secret)
		if err != nil {
			return err
		}

		signedUser := models.SignedUser{
			UserID:   agt.ID.String(),
			UserType: agentType,
			Token:    token,
		}
		_ = ctx.Status(http.StatusOK).JSON(signedUser)

		return nil
	}
}

func RegisterAgent(agentDomain agent.Interactor) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		var params agent.RegistrationParams
		_ = ctx.BodyParser(&params)

		err := params.Validate()
		if err != nil {
			return ctx.Status(http.StatusBadRequest).JSON(err)
		}

		// register agent
		agt, err := agentDomain.Register(params)
		if err != nil {
			return err
		}

		// we use a presenter to reformat the response of agent.
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(agt.ID, models.UserTypAgent))

		return nil
	}
}

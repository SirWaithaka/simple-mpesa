package user_handlers

import (
	"net/http"

	"simple-mpesa/src"
	"simple-mpesa/src/domain/agent"
	"simple-mpesa/src/domain/auth"
	"simple-mpesa/src/routing/responses"
	"simple-mpesa/src/value_objects"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateAgent(agentDomain agent.Interactor, config src.Config) fiber.Handler {

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
		agentType := value_objects.UserTypAgent
		if agt.IsSuperAgent() {
			agentType = value_objects.UserTypSuperAgent
		}

		// generate an auth token string
		token, err := auth.GetTokenString(agt.ID, agentType, config.Secret)
		if err != nil {
			return err
		}

		signedUser := value_objects.SignedUser{
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
		_ = ctx.Status(http.StatusOK).JSON(responses.RegistrationResponse(agt.ID, value_objects.UserTypAgent))

		return nil
	}
}

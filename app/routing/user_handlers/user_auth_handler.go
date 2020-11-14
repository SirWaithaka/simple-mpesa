package user_handlers

import (
	"simple-mpesa/app"
	"simple-mpesa/app/models"
	"simple-mpesa/app/registry"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(domain *registry.Domain, config app.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		// get the user type authenticating
		userType := ctx.Params("user_type")

		switch models.UserType(userType) {
		case models.UserTypAdmin:
			return AuthenticateAdmin(domain.Admin, config)(ctx)
		case models.UserTypAgent:
			return AuthenticateAgent(domain.Agent, config)(ctx)
		case models.UserTypMerchant:
			return AuthenticateMerchant(domain.Merchant, config)(ctx)
		case models.UserTypSubscriber:
			return AuthenticateSubscriber(domain.Subscriber, config)(ctx)
		default:
			return fiber.ErrNotFound
		}
	}
}

func Register(domain *registry.Domain) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		// get the user type authenticating
		userType := ctx.Params("user_type")

		switch models.UserType(userType) {
		case models.UserTypAdmin:
			return RegisterAdmin(domain.Admin)(ctx)
		case models.UserTypAgent:
			return RegisterAgent(domain.Agent)(ctx)
		case models.UserTypMerchant:
			return RegisterMerchant(domain.Merchant)(ctx)
		case models.UserTypSubscriber:
			return RegisterSubscriber(domain.Subscriber)(ctx)
		default:
			return fiber.ErrNotFound
		}
	}
}

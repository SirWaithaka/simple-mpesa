package user_handlers

import (
	"simple-mpesa/src"
	"simple-mpesa/src/value_objects"
	"simple-mpesa/src/registry"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(domain *registry.Domain, config src.Config) fiber.Handler {

	return func(ctx *fiber.Ctx) error {
		// get the user type authenticating
		userType := ctx.Params("user_type")

		switch value_objects.UserType(userType) {
		case value_objects.UserTypAdmin:
			return AuthenticateAdmin(domain.Admin, config)(ctx)
		case value_objects.UserTypAgent:
			return AuthenticateAgent(domain.Agent, config)(ctx)
		case value_objects.UserTypMerchant:
			return AuthenticateMerchant(domain.Merchant, config)(ctx)
		case value_objects.UserTypSubscriber:
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

		switch value_objects.UserType(userType) {
		case value_objects.UserTypAdmin:
			return RegisterAdmin(domain.Admin)(ctx)
		case value_objects.UserTypAgent:
			return RegisterAgent(domain.Agent)(ctx)
		case value_objects.UserTypMerchant:
			return RegisterMerchant(domain.Merchant)(ctx)
		case value_objects.UserTypSubscriber:
			return RegisterSubscriber(domain.Subscriber)(ctx)
		default:
			return fiber.ErrNotFound
		}
	}
}

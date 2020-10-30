package user_handlers

import (
	"simple-wallet/app/registry"

	"github.com/gofiber/fiber/v2"
)

func Authenticate(domain *registry.Domain) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		return nil
	}
}


func Register(domain *registry.Domain) fiber.Handler {

	return func(ctx *fiber.Ctx) error {

		return nil
	}
}
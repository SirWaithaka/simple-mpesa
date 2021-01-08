package routing

import (
	"simple-mpesa/src"
	"simple-mpesa/src/registry"
	"simple-mpesa/src/routing/account_handlers"
	"simple-mpesa/src/routing/error_handlers"
	"simple-mpesa/src/routing/middleware"
	"simple-mpesa/src/routing/transaction_handlers"
	"simple-mpesa/src/routing/user_handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Router(domain *registry.Domain, config src.Config) *fiber.App {

	srv := fiber.New(
		fiber.Config{ErrorHandler: error_handlers.ErrorHandler},
	)

	apiGroup := srv.Group("/api")
	apiGroup.Use(logger.New())

	apiRouteGroup(apiGroup, domain, config)

	return srv
}

func apiRouteGroup(api fiber.Router, domain *registry.Domain, config src.Config) {

	api.Post("/login/:user_type", user_handlers.Authenticate(domain, config))
	api.Post("/user/:user_type", user_handlers.Register(domain))

	// create group at /api/admin
	admin := api.Group("/admin", middleware.AuthByBearerToken(config.Secret))
	admin.Post("/assign-float", user_handlers.AssignFloat(domain.Admin))
	admin.Post("/update-charge", user_handlers.UpdateCharge(domain.Tariff))
	admin.Get("/get-tariff", user_handlers.GetTariff(domain.Tariff))
	admin.Put("/super-agent-status", user_handlers.UpdateSuperAgentStatus(domain.Agent))

	// create group at /api/account
	account := api.Group("/account", middleware.AuthByBearerToken(config.Secret))
	account.Get("/balance", account_handlers.BalanceEnquiry(domain.Account))
	account.Get("/statement", account_handlers.MiniStatement(domain.Account))

	// create group at /api/transaction
	transaction := api.Group("/transaction", middleware.AuthByBearerToken(config.Secret))
	transaction.Post("/deposit", transaction_handlers.Deposit(domain.Transaction))
	transaction.Post("/transfer", transaction_handlers.Transfer(domain.Transaction))
	transaction.Post("/withdraw", transaction_handlers.Withdraw(domain.Transaction))
}

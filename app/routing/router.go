package routing

import (
	"simple-mpesa/app"
	"simple-mpesa/app/registry"
	"simple-mpesa/app/routing/account_handlers"
	"simple-mpesa/app/routing/error_handlers"
	"simple-mpesa/app/routing/middleware"
	"simple-mpesa/app/routing/user_handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Router(domain *registry.Domain, config app.Config) *fiber.App {

	srv := fiber.New(
		fiber.Config{ErrorHandler: error_handlers.ErrorHandler},
	)

	apiGroup := srv.Group("/api")
	apiGroup.Use(logger.New())

	apiRouteGroup(apiGroup, domain, config)

	return srv
}

func apiRouteGroup(api fiber.Router, domain *registry.Domain, config app.Config) {

	api.Post("/login/:user_type", user_handlers.Authenticate(domain, config))
	api.Post("/user/:user_type", user_handlers.Register(domain))

	transaction := api.Group("/transaction", middleware.AuthByBearerToken(config.Secret))
	transaction.Get("/balance", account_handlers.BalanceEnquiry(domain.Account))

	// api.Get("/account/balance", middleware.AuthByBearerToken(config.Secret), account_handlers.BalanceEnquiry(domain.Account))
	// api.Post("/account/deposit", middleware.AuthByBearerToken(config.Secret), account_handlers.Deposit(domain.Account))
	// api.Post("/account/withdraw", middleware.AuthByBearerToken(config.Secret), account_handlers.Withdraw(domain.Account))
	//
	// api.Get("/account/statement", middleware.AuthByBearerToken(config.Secret), account_handlers.MiniStatement(domain.Transaction))
}

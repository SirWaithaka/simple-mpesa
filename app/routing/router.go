package routing

import (
	"simple-wallet/app"
	"simple-wallet/app/registry"
	"simple-wallet/app/routing/account_handlers"
	"simple-wallet/app/routing/error_handlers"
	"simple-wallet/app/routing/middleware"
	"simple-wallet/app/routing/user_handlers"

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

func apiRouteGroup(g fiber.Router, domain *registry.Domain, config app.Config) {

	g.Post("/login", user_handlers.Authenticate(domain.User, config))
	g.Post("/user", user_handlers.Register(domain.User))

	g.Get("/account/balance", middleware.AuthByBearerToken(config.Secret), account_handlers.BalanceEnquiry(domain.Account))
	g.Post("/account/deposit", middleware.AuthByBearerToken(config.Secret), account_handlers.Deposit(domain.Account))
	g.Post("/account/withdrawal", middleware.AuthByBearerToken(config.Secret), account_handlers.Withdraw(domain.Account))
	g.Post("/account/withdraw", middleware.AuthByBearerToken(config.Secret), account_handlers.Withdraw(domain.Account))

	g.Get("/account/statement", middleware.AuthByBearerToken(config.Secret), account_handlers.MiniStatement(domain.Transaction))
}

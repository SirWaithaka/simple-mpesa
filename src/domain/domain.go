package domain

import (
	"simple-mpesa/src"
	"simple-mpesa/src/domain/account"
	"simple-mpesa/src/domain/admin"
	"simple-mpesa/src/domain/agent"
	"simple-mpesa/src/domain/customer"
	"simple-mpesa/src/domain/merchant"
	"simple-mpesa/src/domain/subscriber"
	"simple-mpesa/src/domain/tariff"
	"simple-mpesa/src/domain/transaction"
	"simple-mpesa/src/registry"
	"simple-mpesa/src/repositories/pg"
	"simple-mpesa/src/storage"
)

type Domain struct {
	Admin      admin.Interactor
	Agent      agent.Interactor
	Merchant   merchant.Interactor
	Subscriber subscriber.Interactor

	Account     account.Interactor
	Transaction transaction.Facade
	Tariff      tariff.Manager
}

func NewDomain(config src.Config, database *storage.Database, channels *registry.Channels) *Domain {
	adminRepo := pg.NewAdminRepository(database)
	agentRepo := pg.NewAgentRepository(database)
	merchantRepo := pg.NewMerchantRepository(database)
	subscriberRepo := pg.NewSubscriberRepository(database)

	accRepo := pg.NewAccountRepository(database)
	// txnRepo := repositories.NewTransactionRepository(database)
	statementRepo := pg.NewStatementRepository(database)
	tariffRepo := pg.NewTariffRepository(database)

	// initialize ports and adapters
	ledger := account.NewLedger(statementRepo)
	tariffManager := tariff.NewManager(tariffRepo)
	accountant := account.NewAccountant(accRepo, ledger)
	customerFinder := customer.NewFinder(agentRepo, merchantRepo, subscriberRepo)
	transactor := transaction.NewTransactor(accountant, tariffManager)

	return &Domain{
		Admin:       admin.NewInteractor(config, adminRepo, accountant, customerFinder),
		Agent:       agent.NewInteractor(config, agentRepo, channels.ChannelNewUsers),
		Merchant:    merchant.NewInteractor(config, merchantRepo, channels.ChannelNewUsers),
		Subscriber:  subscriber.NewInteractor(config, subscriberRepo, channels.ChannelNewUsers),
		Account:     account.NewInteractor(accRepo, statementRepo, channels.ChannelNewUsers, channels.ChannelNewTransactions),
		Transaction: transaction.NewFacade(customerFinder, transactor),
		Tariff:      tariffManager,
	}
}

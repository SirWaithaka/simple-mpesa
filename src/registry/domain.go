package registry

import (
	"simple-mpesa/src"
	"simple-mpesa/src/account"
	"simple-mpesa/src/admin"
	"simple-mpesa/src/agent"
	"simple-mpesa/src/customer"
	"simple-mpesa/src/merchant"
	"simple-mpesa/src/ports"
	"simple-mpesa/src/repositories"
	"simple-mpesa/src/statement"
	"simple-mpesa/src/storage"
	"simple-mpesa/src/subscriber"
	"simple-mpesa/src/tariff"
	"simple-mpesa/src/transaction"
)

type Domain struct {
	Admin      admin.Interactor
	Agent      agent.Interactor
	Merchant   merchant.Interactor
	Subscriber subscriber.Interactor

	Account     account.Interactor
	Transaction transaction.Interactor
	Statement   statement.Interactor
	Tariff      tariff.Manager

	Transactor ports.TransactorPort
}

func NewDomain(config src.Config, database *storage.Database, channels *Channels) *Domain {
	adminRepo := repositories.NewAdminRepository(database)
	agentRepo := repositories.NewAgentRepository(database)
	merchantRepo := repositories.NewMerchantRepository(database)
	subscriberRepo := repositories.NewSubscriberRepository(database)

	accRepo := repositories.NewAccountRepository(database)
	txnRepo := repositories.NewTransactionRepository(database)
	statementRepo := repositories.NewStatementRepository(database)
	tariffRepo := repositories.NewTariffRepository(database)

	// initialize ports and adapters
	ledger := statement.NewLedger(statementRepo)
	tariffManager := tariff.NewManager(tariffRepo)
	accountant := account.NewAccountant(accRepo, ledger)
	customerFinder := customer.NewFinder(agentRepo, merchantRepo, subscriberRepo)
	transactor := transaction.NewTransactor(accountant, tariffManager)

	return &Domain{
		Admin:       admin.NewInteractor(config, adminRepo, accountant, customerFinder),
		Agent:       agent.NewInteractor(config, agentRepo, channels.ChannelNewUsers),
		Merchant:    merchant.NewInteractor(config, merchantRepo, channels.ChannelNewUsers),
		Subscriber:  subscriber.NewInteractor(config, subscriberRepo, channels.ChannelNewUsers),
		Account:     account.NewInteractor(accRepo, channels.ChannelNewUsers, channels.ChannelNewTransactions),
		Transaction: transaction.NewInteractor(txnRepo, channels.ChannelNewTransactions),
		Statement:   statement.NewInteractor(statementRepo),
		Transactor:  ports.NewTransactor(customerFinder, transactor),
		Tariff:      tariffManager,
	}
}

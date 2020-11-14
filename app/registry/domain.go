package registry

import (
	"simple-mpesa/app"
	"simple-mpesa/app/account"
	"simple-mpesa/app/admin"
	"simple-mpesa/app/agent"
	"simple-mpesa/app/customer"
	"simple-mpesa/app/merchant"
	"simple-mpesa/app/ports"
	"simple-mpesa/app/statement"
	"simple-mpesa/app/storage"
	"simple-mpesa/app/subscriber"
	"simple-mpesa/app/tariff"
	"simple-mpesa/app/transaction"
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

func NewDomain(config app.Config, database *storage.Database, channels *Channels) *Domain {
	adminRepo := admin.NewRepository(database)
	agentRepo := agent.NewRepository(database)
	merchantRepo := merchant.NewRepository(database)
	subscriberRepo := subscriber.NewRepository(database)

	accRepo := account.NewRepository(database)
	txnRepo := transaction.NewRepository(database)
	statementRepo := statement.NewRepository(database)
	tariffRepo := tariff.NewRepository(database)

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

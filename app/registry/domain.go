package registry

import (
	"simple-wallet/app"
	"simple-wallet/app/account"
	"simple-wallet/app/admin"
	"simple-wallet/app/agent"
	"simple-wallet/app/merchant"
	"simple-wallet/app/storage"
	"simple-wallet/app/subscriber"
	"simple-wallet/app/transaction"
)

type Domain struct {
	Admin      admin.Interactor
	Agent      agent.Interactor
	Merchant   merchant.Interactor
	Subscriber subscriber.Interactor

	Account     account.Interactor
	Transaction transaction.Interactor
}

func NewDomain(config app.Config, database *storage.Database, channels *Channels) *Domain {
	adminRepo := admin.NewRepository(database)
	agentRepo := agent.NewRepository(database)
	merchantRepo := merchant.NewRepository(database)
	subscriberRepo := subscriber.NewRepository(database)

	accRepo := account.NewRepository(database)
	txnRepo := transaction.NewRepository(database)

	return &Domain{
		Admin:       admin.NewInteractor(config, adminRepo),
		Agent:       agent.NewInteractor(config, agentRepo, channels.ChannelNewUsers),
		Merchant:    merchant.NewInteractor(config, merchantRepo, channels.ChannelNewUsers),
		Subscriber:  subscriber.NewInteractor(config, subscriberRepo, channels.ChannelNewUsers),
		Account:     account.NewInteractor(accRepo, channels.ChannelNewUsers, channels.ChannelNewTransactions),
		Transaction: transaction.NewInteractor(txnRepo, channels.ChannelNewTransactions),
	}
}

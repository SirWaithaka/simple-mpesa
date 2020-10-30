package registry

import (
	"simple-wallet/app"
	"simple-wallet/app/account"
	"simple-wallet/app/storage"
	"simple-wallet/app/transaction"
	"simple-wallet/app/user"
)

type Domain struct {
	Account account.Interactor
	Transaction transaction.Interactor
	User user.Interactor
}

func NewDomain(config app.Config, database *storage.Database, channels *Channels) *Domain {
	accRepo := account.NewRepository(database)
	userRepo := user.NewRepository(database)
	txnRepo := transaction.NewRepository(database)

	return &Domain{
		Account: account.NewInteractor(accRepo, channels.ChannelNewUsers, channels.ChannelNewTransactions),
		Transaction: transaction.NewInteractor(txnRepo, channels.ChannelNewTransactions),
		User: user.NewInteractor(config, userRepo, channels.ChannelNewUsers),
	}
}
package registry

import "simple-wallet/app/data"

type Channels struct {
	ChannelNewUsers        data.ChanNewUsers
	ChannelNewTransactions data.ChanNewTransactions
}

func NewChannels() *Channels {
	chanNewUsers := make(chan data.UserContract, 10)
	chanNewTransactions := make(chan data.TransactionContract, 50)

	return &Channels{
		ChannelNewUsers: data.ChanNewUsers{
			Channel: chanNewUsers,
			Reader:  chanNewUsers,
			Writer:  chanNewUsers,
		},
		ChannelNewTransactions: data.ChanNewTransactions{
			Channel: chanNewTransactions,
			Reader:  chanNewTransactions,
			Writer:  chanNewTransactions,
		},
	}
}

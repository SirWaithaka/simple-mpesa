package registry

import (
	"simple-mpesa/app/data"
)

type Channels struct {
	ChannelNewUsers        data.ChanNewCustomers
	ChannelNewTransactions data.ChanNewTransactions
	// ChannelTxnEvents       data.ChanNewTxnEvents
}

func NewChannels() *Channels {
	chanNewUsers := make(chan data.CustomerContract, 10)
	chanNewTransactions := make(chan data.TransactionContract, 50)
	// chanNewTxnEvents := make(chan models.TxnEvent, 100)

	return &Channels{
		ChannelNewUsers: data.ChanNewCustomers{
			Channel: chanNewUsers,
			Reader:  chanNewUsers,
			Writer:  chanNewUsers,
		},
		ChannelNewTransactions: data.ChanNewTransactions{
			Channel: chanNewTransactions,
			Reader:  chanNewTransactions,
			Writer:  chanNewTransactions,
		},
		// ChannelTxnEvents: data.ChanNewTxnEvents{
		// 	Channel: chanNewTxnEvents,
		// 	Reader:  chanNewTxnEvents,
		// 	Writer:  chanNewTxnEvents,
		// },
	}
}

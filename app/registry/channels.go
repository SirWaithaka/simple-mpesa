package registry

import (
	"simple-mpesa/app/data"
)

type Channels struct {
	ChannelNewUsers        data.ChanNewCustomers
}

func NewChannels() *Channels {
	chanNewUsers := make(chan data.CustomerContract, 10)

	return &Channels{
		ChannelNewUsers: data.ChanNewCustomers{
			Channel: chanNewUsers,
			Reader:  chanNewUsers,
			Writer:  chanNewUsers,
		},
	}
}

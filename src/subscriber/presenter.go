package subscriber

import (
	"simple-mpesa/src/data"
)

func parseToNewSubscriber(subscriber Subscriber) data.CustomerContract {
	return data.CustomerContract{
		UserID: subscriber.ID,
	}
}

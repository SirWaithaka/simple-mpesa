package subscriber

import (
	"simple-wallet/app/data"
	"simple-wallet/app/models"
)

func parseToNewSubscriber(subscriber models.Subscriber) data.CustomerContract {
	return data.CustomerContract{
		UserID: subscriber.ID,
	}
}

package subscriber

import (
	"simple-mpesa/src/data"
	"simple-mpesa/src/models"
)

func parseToNewSubscriber(subscriber models.Subscriber) data.CustomerContract {
	return data.CustomerContract{
		UserID: subscriber.ID,
	}
}

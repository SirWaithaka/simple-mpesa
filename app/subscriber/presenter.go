package subscriber

import (
	"simple-mpesa/app/data"
	"simple-mpesa/app/models"
)

func parseToNewSubscriber(subscriber models.Subscriber) data.CustomerContract {
	return data.CustomerContract{
		UserID: subscriber.ID,
	}
}

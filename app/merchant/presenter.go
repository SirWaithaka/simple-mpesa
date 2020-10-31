package merchant

import (
	"simple-mpesa/app/data"
	"simple-mpesa/app/models"
)

func parseToNewMerchant(merchant models.Merchant) data.CustomerContract {
	return data.CustomerContract{
		UserID: merchant.ID,
	}
}

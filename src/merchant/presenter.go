package merchant

import (
	"simple-mpesa/src/data"
	"simple-mpesa/src/models"
)

func parseToNewMerchant(merchant models.Merchant) data.CustomerContract {
	return data.CustomerContract{
		UserID: merchant.ID,
	}
}

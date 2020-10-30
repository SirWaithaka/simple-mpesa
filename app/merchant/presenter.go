package merchant

import (
	"simple-wallet/app/data"
	"simple-wallet/app/models"
)

func parseToNewMerchant(merchant models.Merchant) data.CustomerContract {
	return data.CustomerContract{
		UserID: merchant.ID,
	}
}

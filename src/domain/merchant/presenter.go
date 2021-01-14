package merchant

import (
	"simple-mpesa/src/data"
)

func parseToNewMerchant(merchant Merchant) data.CustomerContract {
	return data.CustomerContract{
		UserID: merchant.ID,
	}
}

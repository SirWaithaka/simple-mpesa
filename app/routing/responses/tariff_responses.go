package responses

import (
	"simple-mpesa/app/models"
	"simple-mpesa/app/tariff"

	"github.com/gofrs/uuid"
)

type tariffResponse struct {
	ID          uuid.UUID           `json:"id"`
	Operation   models.TxnOperation `json:"txnOperation"`
	Source      models.UserType     `json:"srcUserType"`
	Destination models.UserType     `json:"destUserType"`
	Fee         models.Cents        `json:"fee"`
}

func TariffResponse(charges []tariff.Charge) SuccessResponse {
	var tarif []tariffResponse
	for _, charge := range charges {
		tarif = append(tarif, tariffResponse{
			ID:          charge.ID,
			Operation:   charge.Transaction,
			Source:      charge.SourceUserType,
			Destination: charge.DestinationUserType,
			Fee:         charge.Fee,
		})
	}

	msg := "Tariff retrieved"
	return successResponse(msg, tarif)
}

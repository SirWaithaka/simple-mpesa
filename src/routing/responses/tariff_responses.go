package responses

import (
	"simple-mpesa/src/domain/tariff"
	"simple-mpesa/src/domain/value_objects"

	"github.com/gofrs/uuid"
)

type tariffResponse struct {
	ID          uuid.UUID                  `json:"id"`
	Operation   value_objects.TxnOperation `json:"txnOperation"`
	Source      value_objects.UserType     `json:"srcUserType"`
	Destination value_objects.UserType     `json:"destUserType"`
	Fee         value_objects.Cents        `json:"fee"`
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

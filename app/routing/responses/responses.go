package responses

import (
	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func successResponse(message string, data interface{}) SuccessResponse {
	if message == "" {
		message = "Success"
	}
	return SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func RegistrationResponse(userID uuid.UUID, userType models.UserType) interface{} {
	data := map[string]interface{}{
		"userID": userID,
		"userType": userType,
	}
	return successResponse("user created", data)
}

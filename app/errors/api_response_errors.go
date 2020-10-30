package errors

import "net/http"

// ApiErrorResponse properties returned by the REST api for requests that
// are not successful
type ApiErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// InternalServerError
func InternalServerError(message string) ApiErrorResponse {
	return ApiErrorResponse{
		Error:   "internal error",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// UnauthorizedResponse
func UnauthorizedResponse(message string) ApiErrorResponse {
	return ApiErrorResponse{
		Error:   "unauthorized",
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// BadRequestResponse
func BadRequestResponse(message string) ApiErrorResponse {
	return ApiErrorResponse{
		Error:   "could not process request",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

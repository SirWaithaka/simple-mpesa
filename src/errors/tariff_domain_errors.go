package errors

const (
	ErrTariffNotSet     = ERMessage("transaction fee has not been set")
	ErrChargeExists     = ERMessage("tariff already exists")
	ErrChargeNotFound   = ERMessage("charge not found")
	ErrInvalidOperation = ERMessage("transaction operation not supported")
)

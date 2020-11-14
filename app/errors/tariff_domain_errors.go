package errors

const (
	ErrTariffNotSet     = ERMessage("transaction fee has not been set")
	ErrChargeExists     = ERMessage("tariff already exists")
	ErrInvalidOperation = ERMessage("transaction operation not supported")
)

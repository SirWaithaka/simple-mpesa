package errors

const (
	DepositOnlyAtAgent    = ERMessage("deposit can only be done by an agent")
	WithdrawalOnlyAtAgent = ERMessage("withdrawal can only be done by at an agent")
	CustomerCantDeposit   = ERMessage("customer is not allowed to make a deposit")

	TransactionWithSameAccount = ERMessage("operation not allowed: source and destination accounts similar")
)

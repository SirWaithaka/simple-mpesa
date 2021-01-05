package errors

const (
	DepositOnlyAtAgent     = ERMessage("deposit can only be done by an agent")
	WithdrawalOnlyAtAgent  = ERMessage("withdrawal can only be done by at an agent")
	CustomerCantDeposit    = ERMessage("customer is not allowed to make a deposit")
	SuperAgentCantDeposit  = ERMessage("a super agent is only allowed to deposit for other agents")
	SuperAgentCantTransfer = ERMessage("a super agent is not allowed to make/receive transfers")
	SuperAgentCantWithdraw = ERMessage("a super agent is not allowed to withdraw or do withdrawals")

	TransactionWithSameAccount = ERMessage("operation not allowed: source and destination accounts similar")
)

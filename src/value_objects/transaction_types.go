package value_objects

type TxnOperation string

// IsPrimaryOperation returns true if the given operation is among 'withdraw', 'deposit' or 'transfer'
func (operation TxnOperation) IsPrimaryOperation() bool {
	validOps := [3]TxnOperation{TxnOpDeposit, TxnOpWithdraw, TxnOpTransfer}
	for _, op := range validOps {
		if op == operation {
			return true
		}
	}
	return false
}

const (
	TxnOpDeposit  = TxnOperation("DEPOSIT")
	TxnOpWithdraw = TxnOperation("WITHDRAW")
	TxnOpTransfer = TxnOperation("TRANSFER")

	// only used when an admin is assigning float to a super agent
	TxnFloatAssignment = TxnOperation("FLOAT_ASSIGNMENT")
)

package transaction

import "simple-mpesa/app/models"

type Transaction struct {
	Source      models.TxnCustomer // where money is coming from
	Destination models.TxnCustomer // where money is going
	// we can further use this field to describe the specific type of transaction/transfer
	TxnOperation models.TxnOperation
	// amount of money if shillings being transacted
	Amount models.Shillings
}

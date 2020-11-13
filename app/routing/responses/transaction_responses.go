package responses

import (
	"fmt"
	"time"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type transactionStatement struct {
	ID        uuid.UUID           `json:"transactionId"`
	Type      models.TxnOperation `json:"transactionType"`
	Timestamp time.Time           `json:"timestamp"`
	Amount    float64             `json:"amount"`
	UserID    uuid.UUID           `json:"userId"`
	AccountID uuid.UUID           `json:"accountId"`
}

type miniStatementResponse struct {
	Message string    `json:"message"`
	UserID  uuid.UUID `json:"userID"`

	Transactions []transactionStatement `json:"transactions"`
}

func MiniStatementResponse(userID uuid.UUID, transactions []models.Transaction) SuccessResponse {
	var txns []transactionStatement
	for _, txn := range transactions {
		txns = append(txns, transactionStatement{
			ID:        txn.ID,
			Type:      txn.Operation,
			Timestamp: txn.Timestamp,
			Amount:    txn.Amount,
			UserID:    txn.UserID,
			AccountID: txn.AccountID,
		})
	}

	msg := "mini statement retrieved for the past 5 transactions"
	data := miniStatementResponse{
		Message:      msg,
		UserID:       userID,
		Transactions: txns,
	}

	return successResponse(msg, data)
}

type transactionResponse struct {
	Message string    `json:"message"`
}

func TransactionResponse() SuccessResponse {
	data := transactionResponse{
		Message: "Transaction under processing. You will receive a message shortly.",
	}
	return successResponse("", data)
}

type balanceResponse struct {
	UserID  uuid.UUID `json:"userID"`
	Balance float64   `json:"balance"`
}

func BalanceResponse(userID uuid.UUID, balance float64) SuccessResponse {
	msg := fmt.Sprintf("Your current balance is %v", balance)

	data := balanceResponse{
		UserID:  userID,
		Balance: balance,
	}
	return successResponse(msg, data)
}

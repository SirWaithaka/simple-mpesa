package responses

import (
	"fmt"
	"time"

	"simple-mpesa/app/models"

	"github.com/gofrs/uuid"
)

type transactionStatement struct {
	ID        uuid.UUID `json:"transactionId"`
	Type      string    `json:"transactionType"`
	Timestamp time.Time `json:"timestamp"`
	Amount    float64   `json:"amount"`
	UserID    uuid.UUID `json:"userId"`
	AccountID uuid.UUID `json:"accountId"`
}

type miniStatementResponse struct {
	Message string `json:"message"`
	UserID  string `json:"userID"`

	Transactions []transactionStatement `json:"transactions"`
}

func MiniStatementResponse(userID string, transactions []models.Transaction) SuccessResponse {
	var txns []transactionStatement
	for _, txn := range transactions {
		txns = append(txns, transactionStatement{
			ID:        txn.ID,
			Type:      txn.Type,
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
	Message string  `json:"message"`
	UserID  string  `json:"userID"`
	Balance float64 `json:"balance"`
}

func TransactionResponse(txType models.TxType, userID string, balance float64) SuccessResponse {
	var msg string
	switch txType {
	case models.TxTypeDeposit:
		msg = fmt.Sprintf("Amount successfully deposited. New balance %v", balance)
	case models.TxTypeWithdrawal:
		msg = fmt.Sprintf("Amount successfully withdrawn. New balance %v", balance)
	}

	data := transactionResponse{
		Message: msg,
		UserID:  userID,
		Balance: balance,
	}
	return successResponse(msg, data)
}

type balanceResponse struct {
	UserID  string  `json:"userID"`
	Balance float64 `json:"balance"`
}

func BalanceResponse(userID string, balance float64) SuccessResponse {
	msg := fmt.Sprintf("Your current balance is %v", balance)

	data := balanceResponse{
		UserID:  userID,
		Balance: balance,
	}
	return successResponse(msg, data)
}

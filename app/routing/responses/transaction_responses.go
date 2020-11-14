package responses

import (
	"fmt"
	"time"

	"simple-mpesa/app/models"
	"simple-mpesa/app/statement"

	"github.com/gofrs/uuid"
)

type transactionStatement struct {
	ID             uuid.UUID           `json:"transactionId"`
	Type           models.TxnOperation `json:"transactionType"`
	CreatedAt      time.Time           `json:"createdAt"`
	CreditedAmount float64             `json:"creditedAmount"`
	DebitedAmount  float64             `json:"debitedAmount"`
	UserID         uuid.UUID           `json:"userId"`
	AccountID      uuid.UUID           `json:"accountId"`
}

type miniStatementResponse struct {
	Message string    `json:"message"`
	UserID  uuid.UUID `json:"userID"`

	Statements []transactionStatement `json:"transactions"`
}

func MiniStatementResponse(userID uuid.UUID, statements []statement.Statement) SuccessResponse {
	var stmts []transactionStatement
	for _, stmt := range statements {
		stmts = append(stmts, transactionStatement{
			ID:             stmt.ID,
			Type:           stmt.Operation,
			CreatedAt:      stmt.CreatedAt,
			CreditedAmount: stmt.CreditAmount,
			DebitedAmount:  stmt.DebitAmount,
			UserID:         stmt.UserID,
			AccountID:      stmt.AccountID,
		})
	}

	msg := "mini statement retrieved for the past 5 transactions"
	data := miniStatementResponse{
		Message:    msg,
		UserID:     userID,
		Statements: stmts,
	}

	return successResponse(msg, data)
}

type transactionResponse struct {
	Message string `json:"message"`
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

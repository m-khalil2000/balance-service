package models

import (
	"github.com/shopspring/decimal"
)

type User struct {
	ID      uint64          `json:"userId"`
	Balance decimal.Decimal `json:"balance"`
}

type Transaction struct {
	ID         string          `json:"transactionId"`
	UserID     uint64          `json:"userId"`
	State      string          `json:"state"`      // "win" or "lose"
	Amount     decimal.Decimal `json:"amount"`     // always positive, 2 decimal precision
	SourceType string          `json:"sourceType"` // "game", "server", "payment"
}

type TransactionPayload struct {
	State         string `json:"state"`
	Amount        string `json:"amount"`
	TransactionID string `json:"transactionId"`
}

type TransactionResponse struct {
	Message    string `json:"message"`
	OldBalance string `json:"oldBalance"`
	NewBalance string `json:"newBalance"`
}

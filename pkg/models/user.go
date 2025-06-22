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

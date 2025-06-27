package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID      uint64          `json:"userId"`
	Balance decimal.Decimal `json:"balance"`
}

type BalanceResponse struct {
	UserID  uint64 `json:"userId"`
	Balance string `json:"balance"`
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

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int           // Maximum number of open connections to the database
	MaxIdleConns    int           // Maximum number of idle connections to the database
	ConnMaxLifetime time.Duration // Maximum amount of time(s) a connection may be reused
	ConnMaxIdleTime time.Duration // Maximum amount of time(s) a connection may be idle
}

type ServerConfig struct {
	Port string
}

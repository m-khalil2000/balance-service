package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Storage interface {
	ProcessTransaction(ctx context.Context, userID uint64, txID, state, source string, amount decimal.Decimal) error
	GetBalance(ctx context.Context, userID uint64) (decimal.Decimal, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (p *PostgresStorage) ProcessTransaction(ctx context.Context, userID uint64, txID, state, source string, amount decimal.Decimal) error {
	txUUID, err := uuid.Parse(txID)
	if err != nil {
		return fmt.Errorf("invalid transactionId: %w", err)
	}

	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM transactions WHERE id = $1)`, txUUID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("transaction already processed") // Return specific error string
	}

	if state != "win" && state != "lose" {
		return fmt.Errorf("invalid state")
	}

	var currentBalance decimal.Decimal
	err = tx.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	newBalance := currentBalance
	if state == "win" {
		newBalance = newBalance.Add(amount)
	} else {
		newBalance = newBalance.Sub(amount)
		if newBalance.IsNegative() {
			return fmt.Errorf("insufficient balance")
		}
	}

	_, err = tx.ExecContext(ctx, `
		WITH updated_user AS (
			UPDATE users SET balance = $1 WHERE id = $2 RETURNING id
		)
		INSERT INTO transactions (id, user_id, amount, state, source_type)
		SELECT $3, $2, $4, $5, $6
		FROM updated_user
	`, newBalance, userID, txUUID, amount, state, source)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetBalance(ctx context.Context, userID uint64) (decimal.Decimal, error) {
	var balance decimal.Decimal
	// Use context timeout for the query
	err := p.db.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1`, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return decimal.Zero, fmt.Errorf("user not found")
		}
		return decimal.Zero, err
	}
	return balance, nil
}

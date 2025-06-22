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

	tx, err := p.db.BeginTx(ctx, nil)
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
		return nil // Idempotent: already processed
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
	} else if state == "lose" {
		newBalance = newBalance.Sub(amount)
		if newBalance.IsNegative() {
			return fmt.Errorf("insufficient balance")
		}
	} else {
		return fmt.Errorf("invalid state")
	}

	_, err = tx.ExecContext(ctx, `UPDATE users SET balance = $1 WHERE id = $2`, newBalance, userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (id, user_id, amount, state, source_type)
		VALUES ($1, $2, $3, $4, $5)
	`, txUUID, userID, amount, state, source)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (p *PostgresStorage) GetBalance(ctx context.Context, userID uint64) (decimal.Decimal, error) {
	var balance decimal.Decimal
	err := p.db.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1`, userID).Scan(&balance)
	if err != nil {
		return decimal.Zero, err
	}
	return balance, nil
}

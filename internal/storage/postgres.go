package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"balance-service/internal/logger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// Storage defines the interface any storage backend must implement.
type Storage interface {
	GetBalance(ctx context.Context, userID uint64) (decimal.Decimal, error)
	ProcessTransaction(ctx context.Context, userID uint64, txID, state, source string, amount decimal.Decimal) (oldBalance decimal.Decimal, newBalance decimal.Decimal, err error)
}

// PostgresStorage implements Storage using a PostgreSQL database.
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a PostgresStorage backed by the given DB connection.
func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

/*
ProcessTransaction processes a user's transaction:
  - checks for duplicate txID
  - updates user balance atomically
  - inserts the transaction record
*/
func (p *PostgresStorage) ProcessTransaction(ctx context.Context, userID uint64, txID, state, source string, amount decimal.Decimal) (oldBalance decimal.Decimal, newBalance decimal.Decimal, err error) {
	txUUID, err := uuid.Parse(txID)
	if err != nil {
		logger.Log.Warn("invalid transaction ID format",
			zap.String("txID", txID),
			zap.Error(err),
		)
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid transactionId: %w", err)
	}

	// Start a database transaction
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	defer tx.Rollback()

	// Check for duplicate transaction ID
	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM transactions WHERE id = $1)`, txUUID).Scan(&exists)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	if exists {
		logger.Log.Info("transaction already processed",
			zap.String("txID", txID),
			zap.Uint64("userID", userID),
		)
		return decimal.Zero, decimal.Zero, fmt.Errorf("transaction already processed") // Return specific error string
	}

	if state != "win" && state != "lose" {
		logger.Log.Warn("invalid transaction state",
			zap.String("state", state),
		)
		return decimal.Zero, decimal.Zero, fmt.Errorf("invalid state")
	}

	// lock row for update, fetch balance
	err = tx.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&oldBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Warn("user not found during transaction processing",
				zap.Uint64("userID", userID),
			)
			return decimal.Zero, decimal.Zero, fmt.Errorf("user not found")
		}
		return decimal.Zero, decimal.Zero, err
	}
	// calculate new balance
	newBalance = oldBalance
	if state == "win" {
		newBalance = newBalance.Add(amount)
	} else {
		newBalance = newBalance.Sub(amount)
		if newBalance.IsNegative() {
			logger.Log.Debug("insufficient balance for transaction",
				zap.Uint64("userID", userID),
				zap.String("state", state),
				zap.String("amount", amount.String()),
				zap.String("oldBalance", oldBalance.String()),
				zap.String("newBalance", newBalance.String()),
			)
			return oldBalance, newBalance, fmt.Errorf("insufficient balance")
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
		logger.Log.Error("failed to update user balance or insert transaction",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return decimal.Zero, decimal.Zero, err
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("failed to commit transaction",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return decimal.Zero, decimal.Zero, err
	}
	// Log successful transaction
	logger.Log.Debug("transaction processed successfully",
		zap.Uint64("userID", userID),
		zap.String("txID", txID),
		zap.String("state", state),
		zap.String("amount", amount.String()),
		zap.String("newBalance", newBalance.String()),
	)

	return oldBalance, newBalance, nil
}

// GetBalance retrieves the user's current balance.
func (p *PostgresStorage) GetBalance(ctx context.Context, userID uint64) (decimal.Decimal, error) {
	var balance decimal.Decimal

	// Use context timeout for the query
	err := p.db.QueryRowContext(ctx, `SELECT balance FROM users WHERE id = $1`, userID).Scan(&balance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Warn("user not found during balance lookup",
				zap.Uint64("userID", userID),
			)
			return decimal.Zero, fmt.Errorf("user not found")
		}

		logger.Log.Error("failed to fetch user balance",
			zap.Uint64("userID", userID),
			zap.Error(err),
		)
		return decimal.Zero, err
	}

	logger.Log.Info("fetched user balance",
		zap.Uint64("userID", userID),
		zap.String("balance", balance.String()),
	)

	return balance, nil
}

package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"balance-service/config"
	"balance-service/internal/logger"
	"balance-service/pkg/models"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

// initialize database connection

func InitDatabase(cfg *models.Config) (*sql.DB, error) {
	var db *sql.DB
	var err error

	maxRetries := 30
	initialBackoff := 1 * time.Second
	maxBackoff := 30 * time.Second
	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", config.DSN(cfg.DB))
		if err != nil {
			logger.Log.Warn("failed to open DB connection",
				zap.Int("attempt", i+1),
				zap.Error(err),
			)
			time.Sleep(backoff)
			backoff = minDuration(backoff*2, maxBackoff)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Ping the database to verify connection
		err = db.PingContext(ctx)

		// Successful connection
		if err == nil {
			log.Printf("Connected to DB after %d attempts", i+1)
			logger.Log.Info("connected to DB", zap.Int("attempt", i+1))

			// Set connection pool parameters
			db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
			db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
			db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

			return db, nil
		}

		// Log the error and retry
		logger.Log.Warn("DB ping failed",
			zap.Int("attempt", i+1),
			zap.Int("maxRetries", maxRetries),
			zap.Error(err),
		)

		cancel()
		time.Sleep(backoff)
		backoff = minDuration(backoff*2, maxBackoff)
	}
	logger.Log.Fatal("failed to connect to DB after max retries",
		zap.Int("maxRetries", maxRetries),
		zap.Error(err),
	)
	return nil, errors.New("database connection failed after max retries")
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

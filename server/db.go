package server

import (
	"context"
	"database/sql"
	"log"
	"time"

	"balance-service/config"
	"balance-service/internal/logger"
	"balance-service/pkg/models"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

// initialize database connection

func InitDatabase(cfg *models.Config) *sql.DB {
	var db *sql.DB
	var err error

	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", config.DSN(cfg.DB))
		if err != nil {
			logger.Log.Warn("failed to open DB connection",
				zap.Int("attempt", i+1),
				zap.Error(err),
			)
			time.Sleep(retryInterval)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			log.Printf("Connected to DB after %d attempts", i+1)
			logger.Log.Info("connected to DB",
				zap.Int("attempt", i+1),
			)

			db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
			db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
			db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

			return db
		}
		logger.Log.Warn("DB ping failed",
			zap.Int("attempt", i+1),
			zap.Int("maxRetries", maxRetries),
			zap.Error(err),
		)
		time.Sleep(retryInterval)

	}
	logger.Log.Fatal("failed to connect to DB after max retries",
		zap.Int("maxRetries", maxRetries),
		zap.Error(err),
	)
	return nil
}

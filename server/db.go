package server

import (
	"context"
	"database/sql"
	"log"
	"time"

	"balance-service/config"

	_ "github.com/lib/pq"
)

// initialize database connection

func InitDatabase(cfg *config.Config) *sql.DB {
	var db *sql.DB
	var err error

	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", cfg.DB.DSN())
		if err != nil {
			log.Printf("Failed to connect to DB: %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			log.Printf("Connected to DB after %d attempts", i+1)

			db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
			db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
			db.SetConnMaxLifetime(cfg.DB.ConnMaxLifetime)
			db.SetConnMaxIdleTime(cfg.DB.ConnMaxIdleTime)

			return db
		}
		db.Close()
		log.Printf("DB ping failed (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)

	}
	log.Fatal("Failed to connect to DB after %d attempts: %v", maxRetries, err)
	return nil
}

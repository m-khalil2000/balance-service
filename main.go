package main

import (
	"balance-service/config"
	"balance-service/internal/handlers"
	"balance-service/internal/logger"
	"balance-service/internal/storage"
	"balance-service/server"
	"log"

	"go.uber.org/zap"

	"net/http"
	"os"
	"runtime/debug"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	logger.Init()
	defer logger.Log.Sync()

	// Global panic recovery for main()
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Error("UNCAUGHT PANIC in main", zap.Any("error", r), zap.ByteString("stack", debug.Stack()))
			os.Exit(1)
		}
	}()

	cfg := config.LoadFromEnv()

	// Initialize DB connection
	dbConn, err := server.InitDatabase(cfg)
	defer func() {
		if err = dbConn.Close(); err != nil {
			logger.Log.Fatal("Failed to close DB connection", zap.Error(err))
		}
	}()

	logger.Log.Info("Database connected")

	repo := storage.NewPostgresStorage(dbConn)
	h := handlers.NewHandler(repo)

	r := server.SetupRouter(h)

	// Start server
	port := config.GetPort()
	log.Printf("Starting server on port %s\n", port)

	srv := &http.Server{
		Addr:           ":" + config.GetPort(&cfg.Server),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    90 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Fatal(srv.ListenAndServe())
}

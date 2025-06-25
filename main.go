package main

import (
	"balance-service/config"
	"balance-service/internal/handlers"
	"balance-service/internal/storage"
	"balance-service/server"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	log.Printf("GOMAXPROCS set to: %d", runtime.GOMAXPROCS(0))

	// Global panic recovery for main()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("UNCAUGHT PANIC in main: %v\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.LoadFromEnv()

	// Initialize DB connection
	dbConn := server.InitDatabase(cfg)
	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
	}()

	// Log connection pool configuration
	log.Printf("Database connected")

	repo := storage.NewPostgresStorage(dbConn)
	h := handlers.NewHandler(repo)

	r := server.SetupRouter(h)

	// Start server
	port := cfg.Server.GetPort()
	log.Printf("Starting server on port %s\n", port)

	srv := &http.Server{
		Addr:           ":" + cfg.Server.GetPort(),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    90 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Fatal(srv.ListenAndServe())
}

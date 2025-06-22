package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"balance-service/config"
	"balance-service/internal/handlers"
	"balance-service/internal/storage"
)

func main() {
	cfg := config.LoadFromEnv()

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Ensure the database is ready
	repo := storage.NewPostgresStorage(db)
	h := handlers.NewHandler(repo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/user/{userId}/transaction", h.HandleTransaction)
	r.Get("/user/{userId}/balance", h.HandleGetBalance)

	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

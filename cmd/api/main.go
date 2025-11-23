package main

import (
	"context"
	"log"
	"net/http"

	"github.com/OZIOisgood/gamma/internal/api"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	tools.PrintBanner()
	tools.LoadEnv()

	dbURL := tools.GetEnv("DB_URL")
	natsURL := tools.GetEnv("NATS_URL")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	eventBus, err := events.NewEventBus(natsURL)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v\n", err)
	}
	defer eventBus.Close()

	// Ensure JetStream stream exists
	if err := eventBus.EnsureStream("GAMMA", []string{"uploads.>"}); err != nil {
		log.Fatalf("Failed to ensure NATS stream: %v\n", err)
	}

	srv := api.NewServer(pool, eventBus)

	log.Println("Gamma API listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}

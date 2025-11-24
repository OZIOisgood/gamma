package main

import (
	"context"
	"log"
	"net/http"

	"os"

	"github.com/OZIOisgood/gamma/internal/api"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	tools.PrintBanner("assets/api-banner.txt", color.FgCyan)
	tools.LoadEnv()

	dbURL := tools.GetEnv("DB_URL")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	eventBus, err := events.NewEventBus(natsURL)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v", err)
	}
	defer eventBus.Close()

	srv := api.NewServer(pool, eventBus)

	log.Println("Gamma API listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}

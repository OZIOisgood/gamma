package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/OZIOisgood/gamma/internal/worker"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	tools.PrintBanner("assets/worker-banner.txt", color.FgHiMagenta)
	tools.LoadEnv()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	dbURL := tools.GetEnv("DB_URL")

	workerName := os.Getenv("WORKER_NAME")
	if workerName == "" {
		workerName = "worker-1"
	}

	log.Printf("Starting worker: %s", workerName)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := db.New(pool)
	handler := worker.NewHandler(queries, workerName)

	eventBus, err := events.NewEventBus(natsURL)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v", err)
	}
	defer eventBus.Close()

	// Ensure stream exists for MinIO events
	// MinIO publishes to subjects like "gamma.minio.uploaded"
	if err := eventBus.EnsureStream("GAMMA_MINIO", []string{"gamma.minio.>"}); err != nil {
		log.Fatalf("Failed to ensure NATS stream: %v", err)
	}

	// Subscribe to MinIO upload events
	_, err = eventBus.Subscribe("gamma.minio.uploaded", "transcoding-workers", handler.HandleUploadEvent)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	log.Println("Worker listening for events...")

	// Wait for interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Worker shutting down...")
}

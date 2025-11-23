package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/nats-io/nats.go"
)

func main() {
	tools.LoadEnv()

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	workerName := os.Getenv("WORKER_NAME")
	if workerName == "" {
		workerName = "worker-1"
	}

	log.Printf("Starting worker: %s", workerName)

	eventBus, err := events.NewEventBus(natsURL)
	if err != nil {
		log.Fatalf("Unable to connect to NATS: %v", err)
	}
	defer eventBus.Close()

	// Subscribe to uploads.uploaded events
	// We use a queue group "transcoding-workers" so that if we run multiple workers,
	// the work is distributed among them.
	_, err = eventBus.Subscribe("uploads.uploaded", "transcoding-workers", func(msg *nats.Msg) {
		log.Printf("[%s] Received message on %s: %s", workerName, msg.Subject, string(msg.Data))
		msg.Ack()
	})
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

package worker

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	Queries    *db.Queries
	WorkerName string
}

func NewHandler(queries *db.Queries, workerName string) *Handler {
	return &Handler{
		Queries:    queries,
		WorkerName: workerName,
	}
}

func (h *Handler) HandleUploadEvent(msg *nats.Msg) {
	log.Printf("[%s] Received message on %s", h.WorkerName, msg.Subject)

	var event MinioEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		// If it's bad data, maybe we should ack it to remove it from queue, or term it.
		// For now, let's ack it so we don't loop forever on bad data.
		msg.Ack()
		return
	}

	for _, record := range event.Records {
		// We only care about ObjectCreated events
		if !strings.HasPrefix(record.EventName, "s3:ObjectCreated:") {
			continue
		}

		key := record.S3.Object.Key
		decodedKey, err := url.QueryUnescape(key)
		if err != nil {
			log.Printf("Failed to unescape key %s: %v", key, err)
			continue
		}
		log.Printf("Processing upload for key: %s", decodedKey)

		// Update DB status to 'uploaded'
		_, err = h.Queries.UpdateUploadStatusByKey(context.Background(), db.UpdateUploadStatusByKeyParams{
			S3Key:  decodedKey,
			Status: db.UploadStatusUploaded,
		})
		if err != nil {
			log.Printf("Failed to update upload status for key %s: %v", decodedKey, err)
			// For simplicity, we'll log error and continue. Ideally we should handle partial failures.
			continue
		}

		log.Printf("Successfully marked upload as uploaded: %s", decodedKey)

		// Start HLS processing (placeholder)
		log.Printf("Starting HLS processing for %s...", decodedKey)
	}

	msg.Ack()
}

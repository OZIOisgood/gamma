package webhooks

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Queries *db.Queries
}

func NewHandler(queries *db.Queries) *Handler {
	return &Handler{
		Queries: queries,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/webhooks/minio", h.HandleMinioWebhook)
}

type MinioEvent struct {
	Records []struct {
		EventName string `json:"eventName"`
		S3        struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

func (h *Handler) HandleMinioWebhook(w http.ResponseWriter, r *http.Request) {
	var event MinioEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		log.Printf("Failed to decode webhook event: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
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
		log.Printf("Received upload event for key: %s (decoded: %s)", key, decodedKey)

		// Update DB status to 'uploaded'
		_, err = h.Queries.UpdateUploadStatusByKey(r.Context(), db.UpdateUploadStatusByKeyParams{
			S3Key:  decodedKey,
			Status: db.UploadStatusUploaded,
		})
		if err != nil {
			log.Printf("Failed to update upload status for key %s: %v", decodedKey, err)
			continue
		}
		
		log.Printf("Successfully marked upload as uploaded: %s", decodedKey)
	}

	w.WriteHeader(http.StatusOK)
}

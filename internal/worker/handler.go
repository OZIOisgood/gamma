package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	Queries    *db.Queries
	Storage    *storage.Storage
	WorkerName string
}

func NewHandler(queries *db.Queries, storage *storage.Storage, workerName string) *Handler {
	return &Handler{
		Queries:    queries,
		Storage:    storage,
		WorkerName: workerName,
	}
}

func (h *Handler) HandleUploadEvent(msg *nats.Msg) {
	log.Printf("[%s] Received message on %s", h.WorkerName, msg.Subject)

	var event MinioEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Failed to unmarshal event: %v", err)
		msg.Ack()
		return
	}

	for _, record := range event.Records {
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

		if err := h.processVideo(context.Background(), decodedKey); err != nil {
			log.Printf("Failed to process video %s: %v", decodedKey, err)
		}
	}

	msg.Ack()
}

func (h *Handler) processVideo(ctx context.Context, key string) error {
	// key is like "original/<uploadId>.mp4"
	parts := strings.Split(key, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format: %s", key)
	}
	filename := parts[1]
	uploadIDStr := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Update status to processing
	_, err := h.Queries.UpdateUploadStatusByKey(ctx, db.UpdateUploadStatusByKeyParams{
		S3Key:  key,
		Status: db.UploadStatusProcessing,
	})
	if err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "gamma-worker-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Download file
	localInput := filepath.Join(tmpDir, filename)
	if err := h.Storage.DownloadFile(ctx, key, localInput); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	// Generate Asset ID
	assetID := uuid.New()
	hlsDir := filepath.Join(tmpDir, "hls", assetID.String())
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		return fmt.Errorf("failed to create hls dir: %w", err)
	}

	// Run ffmpeg
	masterPlaylist := filepath.Join(hlsDir, "master.m3u8")
	cmd := exec.Command("ffmpeg",
		"-i", localInput,
		"-profile:v", "baseline",
		"-level", "3.0",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		masterPlaylist,
	)
	// Capture output for debugging
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	// Upload HLS files
	err = filepath.Walk(hlsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(filepath.Join(tmpDir, "hls"), path)
		if err != nil {
			return err
		}
		
		// S3 Key: hls/<assetId>/...
		s3Key := filepath.Join("hls", relPath)
		
		contentType := "application/octet-stream"
		if strings.HasSuffix(path, ".m3u8") {
			contentType = "application/vnd.apple.mpegurl"
		} else if strings.HasSuffix(path, ".ts") {
			contentType = "video/mp2t"
		}

		if err := h.Storage.UploadFile(ctx, s3Key, path, contentType); err != nil {
			return fmt.Errorf("failed to upload %s: %w", s3Key, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to upload HLS files: %w", err)
	}

	// Create Asset record
	var pgAssetID pgtype.UUID
	pgAssetID.Scan(assetID.String())
	var pgUploadID pgtype.UUID
	pgUploadID.Scan(uploadIDStr)

	hlsRoot := fmt.Sprintf("hls/%s/master.m3u8", assetID.String())

	_, err = h.Queries.CreateAsset(ctx, db.CreateAssetParams{
		ID:       pgAssetID,
		UploadID: pgUploadID,
		HlsRoot:  hlsRoot,
		Status:   db.AssetStatusReady,
	})
	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}

	// Update Upload status to done
	_, err = h.Queries.UpdateUploadStatusByKey(ctx, db.UpdateUploadStatusByKeyParams{
		S3Key:  key,
		Status: db.UploadStatusReady,
	})
	if err != nil {
		return fmt.Errorf("failed to update upload status to ready: %w", err)
	}

	log.Printf("Successfully processed video %s -> asset %s", key, assetID.String())
	return nil
}

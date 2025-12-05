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
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	Queries    *db.Queries
	Storage    *storage.Storage
	EventBus   *events.EventBus
	WorkerName string
}

func NewHandler(queries *db.Queries, storage *storage.Storage, eventBus *events.EventBus, workerName string) *Handler {
	return &Handler{
		Queries:    queries,
		Storage:    storage,
		EventBus:   eventBus,
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

	// Run ffmpeg with multi-quality support
	// We will generate 3 variants: 1080p, 720p, 480p
	masterPlaylist := "master.m3u8"
	
	// Ensure output directories exist for variants
	cmd := exec.Command("ffmpeg",
		"-i", localInput,
		"-filter_complex", "[0:v]split=3[v1][v2][v3];[v1]scale=w=1920:h=1080:force_original_aspect_ratio=decrease,pad=ceil(iw/2)*2:ceil(ih/2)*2[v1out];[v2]scale=w=1280:h=720:force_original_aspect_ratio=decrease,pad=ceil(iw/2)*2:ceil(ih/2)*2[v2out];[v3]scale=w=854:h=480:force_original_aspect_ratio=decrease,pad=ceil(iw/2)*2:ceil(ih/2)*2[v3out]",
		
		// 1080p
		"-map", "[v1out]", "-c:v:0", "libx264", "-b:v:0", "5000k", "-maxrate:v:0", "5350k", "-bufsize:v:0", "7500k",
		"-map", "a:0", "-c:a:0", "aac", "-b:a:0", "192k", "-ac", "2",
		
		// 720p
		"-map", "[v2out]", "-c:v:1", "libx264", "-b:v:1", "2800k", "-maxrate:v:1", "2996k", "-bufsize:v:1", "4200k",
		"-map", "a:0", "-c:a:1", "aac", "-b:a:1", "128k", "-ac", "2",
		
		// 480p
		"-map", "[v3out]", "-c:v:2", "libx264", "-b:v:2", "1400k", "-maxrate:v:2", "1498k", "-bufsize:v:2", "2100k",
		"-map", "a:0", "-c:a:2", "aac", "-b:a:2", "96k", "-ac", "2",

		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-master_pl_name", masterPlaylist,
		"-hls_segment_filename", filepath.Join(hlsDir, "v%v_segment%03d.ts"),
		"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2",
		filepath.Join(hlsDir, "v%v.m3u8"),
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

	// Publish asset processed event
	eventData := map[string]string{
		"asset_id":  assetID.String(),
		"upload_id": uploadIDStr,
		"status":    string(db.AssetStatusReady),
	}
	eventBytes, _ := json.Marshal(eventData)
	if err := h.EventBus.Publish("gamma.assets.processed", eventBytes); err != nil {
		log.Printf("Failed to publish asset processed event: %v", err)
	}

	log.Printf("Successfully processed video %s -> asset %s", key, assetID.String())
	return nil
}

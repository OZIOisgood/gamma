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
	"strconv"
	"strings"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nats-io/nats.go"
)

type ProbeData struct {
	Streams []struct {
		Width    int    `json:"width"`
		Height   int    `json:"height"`
		Duration string `json:"duration"`
	} `json:"streams"`
}

type Quality struct {
	Height       int
	Width        int
	Bitrate      string
	MaxRate      string
	BufSize      string
	AudioBitrate string
}

var availableQualities = []Quality{
	{1080, 1920, "5000k", "5350k", "7500k", "192k"},
	{720, 1280, "2800k", "2996k", "4200k", "128k"},
	{480, 854, "1400k", "1498k", "2100k", "96k"},
	{360, 640, "800k", "856k", "1200k", "96k"},
	{240, 426, "400k", "428k", "600k", "64k"},
	{144, 256, "200k", "214k", "300k", "64k"},
}

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

func (h *Handler) getVideoMetadata(filePath string) (int, int, float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration",
		"-of", "json",
		filePath,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0, err
	}

	var data ProbeData
	if err := json.Unmarshal(output, &data); err != nil {
		return 0, 0, 0, err
	}

	if len(data.Streams) == 0 {
		return 0, 0, 0, fmt.Errorf("no video streams found")
	}

	duration, err := strconv.ParseFloat(data.Streams[0].Duration, 64)
	if err != nil {
		// Fallback if duration is missing or invalid, though unlikely for valid video
		duration = 0
	}

	return data.Streams[0].Width, data.Streams[0].Height, duration, nil
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

func (h *Handler) generateThumbnails(ctx context.Context, inputPath, outputDir string, duration float64) (string, error) {
	thumbDir := filepath.Join(outputDir, "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create thumbnails dir: %w", err)
	}

	// Generate thumbnails every 10 seconds
	// Scale to width 160, keep aspect ratio
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", "fps=1/10,scale=160:-1",
		filepath.Join(thumbDir, "thumb%03d.jpg"),
	)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg thumbnails failed: %w", err)
	}

	// Generate VTT file
	vttPath := filepath.Join(thumbDir, "thumbnails.vtt")
	f, err := os.Create(vttPath)
	if err != nil {
		return "", fmt.Errorf("failed to create vtt file: %w", err)
	}
	defer f.Close()

	f.WriteString("WEBVTT\n\n")

	// List generated files
	files, err := filepath.Glob(filepath.Join(thumbDir, "thumb*.jpg"))
	if err != nil {
		return "", fmt.Errorf("failed to list thumbnails: %w", err)
	}

	for i := range files {
		start := i * 10
		end := (i + 1) * 10

		// Format time as HH:MM:SS.mmm
		startTime := fmt.Sprintf("%02d:%02d:%02d.000", start/3600, (start%3600)/60, start%60)
		endTime := fmt.Sprintf("%02d:%02d:%02d.000", end/3600, (end%3600)/60, end%60)

		f.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
		f.WriteString(fmt.Sprintf("thumb%03d.jpg\n\n", i+1))
	}

	return "thumbnails/thumbnails.vtt", nil
}

func (h *Handler) processVideo(ctx context.Context, key string) error {
	// key format: realm/original/<uploadId>.mp4
	parts := strings.Split(key, "/")

	if len(parts) != 3 || parts[1] != "original" {
		// Not an original file upload, ignore (could be HLS segments etc)
		return nil
	}

	realmPrefix := parts[0]
	filename := parts[2]
	uploadIDStr := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Update status to processing and get the upload (to access realm_id)
	upload, err := h.Queries.UpdateUploadStatusByKey(ctx, db.UpdateUploadStatusByKeyParams{
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

	// Get video metadata
	width, height, duration, err := h.getVideoMetadata(localInput)
	if err != nil {
		return fmt.Errorf("failed to get video metadata: %w", err)
	}
	log.Printf("Input video dimensions: %dx%d, duration: %f", width, height, duration)

	// Filter qualities
	var targetQualities []Quality
	for _, q := range availableQualities {
		if q.Height <= height {
			targetQualities = append(targetQualities, q)
		}
	}
	if len(targetQualities) == 0 {
		// Fallback to lowest quality if input is smaller than 144p
		targetQualities = append(targetQualities, availableQualities[len(availableQualities)-1])
	}

	// Generate Asset ID
	assetID := uuid.New()
	hlsDir := filepath.Join(tmpDir, "hls", assetID.String())
	if err := os.MkdirAll(hlsDir, 0755); err != nil {
		return fmt.Errorf("failed to create hls dir: %w", err)
	}

	// Generate Thumbnails
	vttRelPath, err := h.generateThumbnails(ctx, localInput, hlsDir, duration)
	if err != nil {
		log.Printf("Failed to generate thumbnails: %v", err)
	}

	// Run ffmpeg with multi-quality support
	masterPlaylist := "master.m3u8"

	// Build ffmpeg command dynamically
	args := []string{"-i", localInput}

	// Filter complex
	filterComplex := fmt.Sprintf("[0:v]split=%d", len(targetQualities))
	for i := 0; i < len(targetQualities); i++ {
		filterComplex += fmt.Sprintf("[v%d]", i+1)
	}
	filterComplex += ";"
	for i, q := range targetQualities {
		filterComplex += fmt.Sprintf("[v%d]scale=w=%d:h=%d:force_original_aspect_ratio=decrease,pad=ceil(iw/2)*2:ceil(ih/2)*2[v%dout];", i+1, q.Width, q.Height, i+1)
	}
	filterComplex = strings.TrimSuffix(filterComplex, ";")
	args = append(args, "-filter_complex", filterComplex)

	// Map streams
	varStreamMap := ""
	for i, q := range targetQualities {
		args = append(args,
			"-map", fmt.Sprintf("[v%dout]", i+1),
			fmt.Sprintf("-c:v:%d", i), "libx264",
			fmt.Sprintf("-b:v:%d", i), q.Bitrate,
			fmt.Sprintf("-maxrate:v:%d", i), q.MaxRate,
			fmt.Sprintf("-bufsize:v:%d", i), q.BufSize,
			"-map", "a:0",
			fmt.Sprintf("-c:a:%d", i), "aac",
			fmt.Sprintf("-b:a:%d", i), q.AudioBitrate,
			"-ac", "2",
		)
		varStreamMap += fmt.Sprintf("v:%d,a:%d ", i, i)
	}

	args = append(args,
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-master_pl_name", masterPlaylist,
		"-hls_segment_filename", filepath.Join(hlsDir, "v%v_segment%03d.ts"),
		"-var_stream_map", strings.TrimSpace(varStreamMap),
		filepath.Join(hlsDir, "v%v.m3u8"),
	)

	cmd := exec.Command("ffmpeg", args...)
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

		// S3 Key: [realm/]hls/<assetId>/...
		s3Key := filepath.Join(realmPrefix, "hls", relPath)

		contentType := "application/octet-stream"
		if strings.HasSuffix(path, ".m3u8") {
			contentType = "application/vnd.apple.mpegurl"
		} else if strings.HasSuffix(path, ".ts") {
			contentType = "video/mp2t"
		} else if strings.HasSuffix(path, ".vtt") {
			contentType = "text/vtt"
		} else if strings.HasSuffix(path, ".jpg") {
			contentType = "image/jpeg"
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

	hlsRoot := filepath.Join(realmPrefix, "hls", assetID.String(), "master.m3u8")

	var thumbnailRoot pgtype.Text
	if vttRelPath != "" {
		thumbnailRoot.Scan(filepath.Join(realmPrefix, "hls", assetID.String(), vttRelPath))
	}

	_, err = h.Queries.CreateAsset(ctx, db.CreateAssetParams{
		ID:            pgAssetID,
		UploadID:      pgUploadID,
		RealmID:       upload.RealmID,
		HlsRoot:       hlsRoot,
		ThumbnailRoot: thumbnailRoot,
		Status:        db.AssetStatusReady,
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

func (h *Handler) HandleDeleteAssetEvent(msg *nats.Msg) {
	log.Printf("[%s] Received delete asset message", h.WorkerName)

	var event map[string]string
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("Failed to unmarshal delete event: %v", err)
		msg.Ack()
		return
	}

	assetID := event["asset_id"]
	uploadID := event["upload_id"]
	hlsRoot := event["hls_root"]

	ctx := context.Background()

	// Delete original file
	if uploadID != "" {
		var pgUUID pgtype.UUID
		pgUUID.Scan(uploadID)
		upload, err := h.Queries.GetUpload(ctx, pgUUID)
		if err == nil {
			if err := h.Storage.DeleteFile(ctx, upload.S3Key); err != nil {
				log.Printf("Failed to delete original file %s: %v", upload.S3Key, err)
			}
		} else {
			log.Printf("Failed to get upload %s: %v", uploadID, err)
		}
	}

	// Delete HLS folder using hls_root path (e.g., "realm/hls/assetId/master.m3u8")
	if hlsRoot != "" {
		// Get folder path by removing the filename
		hlsFolder := filepath.Dir(hlsRoot) + "/"
		if err := h.Storage.DeleteFolder(ctx, hlsFolder); err != nil {
			log.Printf("Failed to delete HLS folder %s: %v", hlsFolder, err)
		}
	} else if assetID != "" {
		// Fallback: try without realm prefix (shouldn't happen in MVP)
		log.Printf("Warning: hls_root not provided, using fallback path for asset %s", assetID)
	}

	msg.Ack()
}

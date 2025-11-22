package uploads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	Storage *storage.Storage
	Queries *db.Queries
}

func NewHandler(storage *storage.Storage, queries *db.Queries) *Handler {
	return &Handler{
		Storage: storage,
		Queries: queries,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/uploads", h.CreateUpload)
	r.Get("/uploads", h.List)
	r.Get("/uploads/{id}", h.Get)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	videos, err := h.Queries.ListUploads(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list videos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	video, err := h.Queries.GetUpload(r.Context(), pgUUID)
	if err != nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

type CreateUploadRequest struct {
	Filename string `json:"filename"`
}

type CreateUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
}

func (h *Handler) CreateUpload(w http.ResponseWriter, r *http.Request) {
	var req CreateUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Filename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	// Generate a unique ID for the video
	videoID := uuid.New()
	ext := filepath.Ext(req.Filename)
	key := fmt.Sprintf("raw/%s%s", videoID.String(), ext)

	// Generate presigned URL
	ctx := r.Context()
	uploadURL, err := h.Storage.GetPresignedURL(ctx, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate upload URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Save to database
	var pgUUID pgtype.UUID
	err = pgUUID.Scan(videoID.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse UUID: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = h.Queries.CreateUpload(ctx, db.CreateUploadParams{
		ID:     pgUUID,
		Title:  req.Filename,
		S3Key:  key,
		Status: db.UploadStatusPending,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create upload record: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateUploadResponse{
		ID:        videoID.String(),
		UploadURL: uploadURL,
		Key:       key,
	})
}

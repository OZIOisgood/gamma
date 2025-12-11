package uploads

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/OZIOisgood/gamma/internal/events"
	"github.com/OZIOisgood/gamma/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Handler struct {
	Storage  *storage.Storage
	Queries  *db.Queries
	EventBus *events.EventBus
}

func NewHandler(storage *storage.Storage, queries *db.Queries, eventBus *events.EventBus) *Handler {
	return &Handler{
		Storage:  storage,
		Queries:  queries,
		EventBus: eventBus,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/uploads", h.CreateUpload)
	r.Get("/uploads", h.List)
	r.Get("/uploads/{id}", h.Get)
	r.Get("/assets", h.ListAssets)
	r.Get("/assets/{id}", h.GetAsset)
	r.Delete("/assets/{id}", h.DeleteAsset)
	r.Get("/assets/{id}/playlist", h.GetAssetPlaylist)
}

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	realmName := chi.URLParam(r, "realm")
	if realmName == "" {
		http.Error(w, "Realm is required", http.StatusBadRequest)
		return
	}

	realm, err := h.Queries.GetRealmByName(r.Context(), realmName)
	if err != nil {
		http.Error(w, "Realm not found", http.StatusNotFound)
		return
	}

	assets, err := h.Queries.ListAssets(r.Context(), realm.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list assets: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assets)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	realmName := chi.URLParam(r, "realm")
	if realmName == "" {
		http.Error(w, "Realm is required", http.StatusBadRequest)
		return
	}

	realm, err := h.Queries.GetRealmByName(r.Context(), realmName)
	if err != nil {
		http.Error(w, "Realm not found", http.StatusNotFound)
		return
	}

	videos, err := h.Queries.ListUploads(r.Context(), realm.ID)
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
	realmName := chi.URLParam(r, "realm")
	if realmName == "" {
		http.Error(w, "Realm is required", http.StatusBadRequest)
		return
	}

	realm, err := h.Queries.GetRealmByName(r.Context(), realmName)
	if err != nil {
		http.Error(w, "Realm not found", http.StatusNotFound)
		return
	}

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
	if ext == "" {
		ext = ".mp4"
	}
	// Key format: realm/original/uuid.ext
	key := fmt.Sprintf("%s/original/%s%s", realmName, videoID.String(), ext)

	// Generate presigned URL
	ctx := r.Context()
	uploadURL, err := h.Storage.GeneratePresignedPutURL(ctx, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate upload URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Save to database
	var pgUUID pgtype.UUID
	pgUUID.Scan(videoID.String())

	_, err = h.Queries.CreateUpload(ctx, db.CreateUploadParams{
		ID:      pgUUID,
		Title:   req.Filename,
		S3Key:   key,
		Status:  db.UploadStatusPending,
		RealmID: realm.ID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create upload record: %v", err), http.StatusInternalServerError)
		return
	}

	resp := CreateUploadResponse{
		ID:        videoID.String(),
		UploadURL: uploadURL,
		Key:       key,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	asset, err := h.Queries.GetAsset(r.Context(), pgUUID)
	if err != nil {
		// Try to find by upload ID as well, just in case user passed upload ID
		asset, err = h.Queries.GetAssetByUploadID(r.Context(), pgUUID)
		if err != nil {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asset)
}

type GetAssetPlaylistResponse struct {
	URL string `json:"url"`
}

func (h *Handler) GetAssetPlaylist(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	asset, err := h.Queries.GetAsset(r.Context(), pgUUID)
	if err != nil {
		// Try to find by upload ID as well, just in case user passed upload ID
		asset, err = h.Queries.GetAssetByUploadID(r.Context(), pgUUID)
		if err != nil {
			http.Error(w, "Asset not found", http.StatusNotFound)
			return
		}
	}

	if asset.Status != db.AssetStatusReady {
		http.Error(w, "Asset is not ready", http.StatusBadRequest)
		return
	}

	key := asset.HlsRoot
	presignedURL, err := h.Storage.GeneratePresignedGetURL(r.Context(), key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate playlist URL: %v", err), http.StatusInternalServerError)
		return
	}

	// Since the bucket is public for HLS, we can strip the query parameters
	// to provide a cleaner URL and avoid expiration issues.
	u, err := url.Parse(presignedURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse URL: %v", err), http.StatusInternalServerError)
		return
	}
	u.RawQuery = ""
	finalURL := u.String()

	resp := GetAssetPlaylistResponse{
		URL: finalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(idStr)
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	asset, err := h.Queries.SoftDeleteAsset(r.Context(), pgUUID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete asset: %v", err), http.StatusInternalServerError)
		return
	}

	// Also soft delete the upload since the dashboard lists uploads
	if _, err := h.Queries.SoftDeleteUpload(r.Context(), asset.UploadID); err != nil {
		log.Printf("Failed to soft delete upload %v: %v", asset.UploadID, err)
	}

	uploadUUID, err := uuid.FromBytes(asset.UploadID.Bytes[:])
	uploadIDStr := ""
	if err == nil {
		uploadIDStr = uploadUUID.String()
	}

	eventData := map[string]string{
		"asset_id":  idStr,
		"upload_id": uploadIDStr,
		"hls_root":  asset.HlsRoot,
	}

	payload, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("Failed to marshal delete event: %v", err)
	} else {
		if err := h.EventBus.Publish("delete_asset", payload); err != nil {
			log.Printf("Failed to publish delete event: %v", err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

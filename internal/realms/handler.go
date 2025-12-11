package realms

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/OZIOisgood/gamma/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
	r.Get("/realms", h.List)
	r.Post("/realms", h.Create)
	r.Get("/realms/{id}", h.Get)
	r.Delete("/realms/{id}", h.Delete)
}

type RealmResponse struct {
	ID        pgtype.UUID        `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
}

func toRealmResponse(r db.Realm) RealmResponse {
	return RealmResponse{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	realms, err := h.Queries.ListRealms(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list realms: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]RealmResponse, len(realms))
	for i, realm := range realms {
		response[i] = toRealmResponse(realm)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type CreateRealmRequest struct {
	Name string `json:"name"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRealmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	realm, err := h.Queries.CreateRealm(r.Context(), req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create realm: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toRealmResponse(realm))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(idStr); err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	realm, err := h.Queries.GetRealm(r.Context(), pgUUID)
	if err != nil {
		http.Error(w, "Realm not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toRealmResponse(realm))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	var pgUUID pgtype.UUID
	if err := pgUUID.Scan(idStr); err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	// Check if it's the last realm
	realms, err := h.Queries.ListRealms(r.Context())
	if err != nil {
		http.Error(w, "Failed to check realms count", http.StatusInternalServerError)
		return
	}

	if len(realms) <= 1 {
		http.Error(w, "Cannot delete the last realm", http.StatusBadRequest)
		return
	}

	// Soft delete all assets for this realm
	if err := h.Queries.SoftDeleteAssetsByRealmID(r.Context(), pgUUID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete assets: %v", err), http.StatusInternalServerError)
		return
	}

	// Soft delete all uploads for this realm
	if err := h.Queries.SoftDeleteUploadsByRealmID(r.Context(), pgUUID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete uploads: %v", err), http.StatusInternalServerError)
		return
	}

	// Soft delete the realm
	if err := h.Queries.DeleteRealm(r.Context(), pgUUID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete realm: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

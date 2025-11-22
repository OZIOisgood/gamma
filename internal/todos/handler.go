package todos

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	r.Get("/todos", h.listTodos)
	r.Post("/todos", h.createTodo)
	r.Get("/todos/{id}", h.getTodo)
	r.Put("/todos/{id}", h.updateTodo)
	r.Delete("/todos/{id}", h.deleteTodo)
}

func (h *Handler) listTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.Queries.ListTodos(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

type CreateTodoRequest struct {
	Task string `json:"task"`
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.Queries.CreateTodo(r.Context(), db.CreateTodoParams{
		Task:      req.Task,
		Completed: false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todo, err := h.Queries.GetTodo(r.Context(), id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

type UpdateTodoRequest struct {
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.Queries.UpdateTodo(r.Context(), db.UpdateTodoParams{
		ID:        id,
		Task:      req.Task,
		Completed: req.Completed,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.Queries.DeleteTodo(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

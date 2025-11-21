package main

import (
  "log"
  "net/http"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := chi.NewRouter()

  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok"}`))
  })

  log.Println("Gamma API listening on :8080")
  if err := http.ListenAndServe(":8080", r); err != nil {
    log.Fatal(err)
  }
}

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/OZIOisgood/gamma/internal/api"
	"github.com/OZIOisgood/gamma/internal/tools"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	tools.PrintBanner()
	tools.LoadEnv()

	dbURL := tools.GetEnv("DB_URL")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	srv := api.NewServer(pool)

	log.Println("Gamma API listening on :8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatal(err)
	}
}

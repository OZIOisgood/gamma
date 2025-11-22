package tools

import (
	"log"
	"os"
)

func GetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s not set", key)
	}
	return v
}

package tools

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func PrintBanner() {
	data, err := os.ReadFile("assets/banner.txt")
	if err != nil {
		log.Printf("Failed to load banner: %v", err)
		return
	}
	color.Cyan(string(data))
}

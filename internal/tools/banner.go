package tools

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func PrintBanner(path string, attr color.Attribute) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to load banner: %v", err)
		return
	}
	color.New(attr).Println(string(data))
}

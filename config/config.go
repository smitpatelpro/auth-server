package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var once sync.Once

// Config func to get env value
func Config(key string) string {
	once.Do(func() {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}
	})
	return os.Getenv(key)
}

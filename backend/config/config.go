package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath    string // SQLite database file path
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DBPath:    getEnv("DB_PATH", "warehouse.db"), // Default SQLite database file
		JWTSecret: getEnv("JWT_SECRET", "default-secret-key"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


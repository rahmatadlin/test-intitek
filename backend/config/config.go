package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath    string // SQLite database file path
	JWTSecret string
	Port      string
}

func LoadConfig() *Config {
	// Try to load .env file from multiple locations
	// 1. Current directory
	// 2. Executable directory (for production build)
	// 3. Parent directories (for production build in build/bin)
	envPaths := []string{
		".env",
		"../.env",
		"../../.env",
		"../../../.env",
	}

	// Also try relative to executable
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		envPaths = append(envPaths,
			filepath.Join(exeDir, ".env"),
			filepath.Join(exeDir, "..", ".env"),
			filepath.Join(exeDir, "..", "..", ".env"),
			filepath.Join(exeDir, "..", "..", "..", ".env"),
		)
	}

	// Try each path
	loaded := false
	for _, envPath := range envPaths {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env file from: %s", envPath)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("No .env file found in any location, using environment variables")
		log.Println("Tried paths:", envPaths)
	}

	// Get DB path - try multiple locations for production build
	dbPath := getEnv("DB_PATH", "")
	if dbPath == "" {
		// Try to determine path relative to executable
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			// If running from build/bin, use data in build directory
			if filepath.Base(exeDir) == "bin" {
				dbPath = filepath.Join(exeDir, "..", "data", "warehouse.db")
			} else {
				dbPath = filepath.Join(exeDir, "data", "warehouse.db")
			}
		} else {
			// Fallback to current directory
			dbPath = "data/warehouse.db"
		}
	}

	return &Config{
		DBPath:    dbPath,
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

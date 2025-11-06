package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	Port       string
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

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "warehouse_db"),
		JWTSecret:  getEnv("JWT_SECRET", "default-secret-key"),
		Port:       getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}


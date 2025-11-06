package main

import (
	"log"
	"os"
	"path/filepath"
	"warehouse-management/config"
	"warehouse-management/database"
	"warehouse-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Check if running as Wails app
	// Check multiple ways to detect Wails mode:
	// 1. Check for wails.json in current directory or executable directory
	// 2. Check for environment variable set by Wails
	// 3. Check if executable is in build/bin directory (production build)
	isWailsMode := false

	// Check wails.json in current directory
	if _, err := os.Stat("wails.json"); err == nil {
		isWailsMode = true
	} else {
		// Check in executable directory
		if exePath, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exePath)
			if _, err := os.Stat(filepath.Join(exeDir, "wails.json")); err == nil {
				isWailsMode = true
			} else if _, err := os.Stat(filepath.Join(exeDir, "..", "wails.json")); err == nil {
				isWailsMode = true
			}
		}
	}

	// Check environment variable
	if os.Getenv("WAILS_MODE") != "" {
		isWailsMode = true
	}

	// Check if running from build directory (production executable)
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		if filepath.Base(exeDir) == "bin" || filepath.Base(exeDir) == "build" {
			isWailsMode = true
		}
	}

	if isWailsMode {
		// Run as Wails desktop app
		RunWailsApp()
		return
	}

	// Otherwise, run as normal web server
	runWebServer()
}

func runWebServer() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(router, cfg)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Warehouse Management API is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

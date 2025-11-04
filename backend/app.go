package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"warehouse-management/config"
	"warehouse-management/database"
	"warehouse-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Untuk production build, kita akan menggunakan embed dari folder yang sudah di-copy
// Atau kita bisa menggunakan approach yang berbeda - load dari file system
// Tapi karena embed tidak bisa pakai .., kita buat folder di dalam backend
// Untuk sementara, kita gunakan empty embed dan akan load dari file system
var assets embed.FS

// App struct
type App struct {
	ctx    context.Context
	router *gin.Engine
	cfg    *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware - allow all origins for Wails
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
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

	return &App{
		router: router,
		cfg:    cfg,
	}
}

// OnStartup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Start Gin server in background on localhost
	go func() {
		port := a.cfg.Port
		if port == "" {
			port = "8080"
		}
		log.Printf("API Server starting on port %s...", port)
		if err := a.router.Run(":" + port); err != nil {
			log.Printf("Failed to start API server: %v", err)
		}
	}()
}

// OnDomReady is called when the frontend is ready
func (a *App) OnDomReady(ctx context.Context) {
	// Frontend is ready
}

// OnBeforeClose is called when the app is about to quit
func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// OnShutdown is called when the app is shutting down
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup if needed
}

// RunWailsApp runs the Wails application
func RunWailsApp() {
	// Create an instance of the app structure
	app := NewApp()

	// Resolve frontend dist path
	// Try to find frontend/dist relative to backend folder
	frontendDistPath := filepath.Join("..", "frontend", "dist")
	if _, err := os.Stat(frontendDistPath); os.IsNotExist(err) {
		// Try absolute path from current working directory
		wd, err := os.Getwd()
		if err == nil {
			frontendDistPath = filepath.Join(wd, "..", "frontend", "dist")
		}
	}

	// Check if path exists, if not use empty embed (will use dev mode)
	assetServerOptions := &assetserver.Options{}
	if _, err := os.Stat(frontendDistPath); err == nil {
		// Use file system path for development
		assetServerOptions.Assets = os.DirFS(frontendDistPath)
	} else {
		// Use embed for production (will be empty if not embedded)
		assetServerOptions.Assets = assets
		log.Printf("Warning: Frontend dist not found at %s, using embedded assets", frontendDistPath)
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "Warehouse Management System",
		Width:            1280,
		Height:           800,
		AssetServer:      assetServerOptions,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.OnStartup,
		OnDomReady:       app.OnDomReady,
		OnBeforeClose:    app.OnBeforeClose,
		OnShutdown:       app.OnShutdown,
	})

	if err != nil {
		log.Fatal("Error:", err)
	}
}

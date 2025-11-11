package main

import (
	"context"
	"embed"
	"log"
	"os"
	"path/filepath"
	"warehouse-management/config"
	"warehouse-management/database"
	"warehouse-management/logger"
	"warehouse-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Embed frontend dist dari folder frontend-dist (harus di-copy sebelum build)
// Untuk build production:
// 1. Copy folder frontend/dist ke backend/frontend-dist
// 2. Build dengan wails build
//
//go:embed all:frontend-dist
var assets embed.FS

// App struct
type App struct {
	ctx        context.Context
	router     *gin.Engine
	cfg        *config.Config
	fileLogger *logger.FileLogger // Wails file logger
}

// GetLogger returns the file logger instance
func (a *App) GetLogger() *logger.FileLogger {
	return a.fileLogger
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	// For Wails app, we need database to work properly
	if err := database.ConnectDatabase(cfg); err != nil {
		log.Printf("ERROR: Failed to connect to database: %v", err)
		log.Printf("Please check:")
		log.Printf("  1. Database directory is writable")
		log.Printf("  2. Database path: %s", cfg.DBPath)
		log.Printf("  3. Check file permissions")
		// Return error but don't fatal - let Wails start and show error in UI
		// Database operations will fail but at least app starts
		log.Printf("Application will start but database features will not work")
	} else {
		log.Printf("Database connected successfully")
		log.Printf("Database path: %s", cfg.DBPath)
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
		if a.fileLogger != nil {
			a.fileLogger.Info("API Server starting on http://localhost:" + port)
			a.fileLogger.Info("API endpoints available at http://localhost:" + port + "/api")
		}
		log.Printf("API Server starting on http://localhost:%s...", port)
		log.Printf("API endpoints available at http://localhost:%s/api", port)
		if err := a.router.Run(":" + port); err != nil {
			if a.fileLogger != nil {
				a.fileLogger.Error("Failed to start API server: " + err.Error())
			}
			log.Printf("Failed to start API server: %v", err)
		}
	}()
}

// OnDomReady is called when the frontend is ready
func (a *App) OnDomReady(ctx context.Context) {
	// Frontend is ready
	if a.fileLogger != nil {
		a.fileLogger.Info("Frontend DOM ready")
	}
}

// OnBeforeClose is called when the app is about to quit
func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	if a.fileLogger != nil {
		a.fileLogger.Info("Application closing...")
	}
	return false
}

// OnShutdown is called when the app is shutting down
func (a *App) OnShutdown(ctx context.Context) {
	// Cleanup if needed
	if a.fileLogger != nil {
		a.fileLogger.Info("Application shutdown")
	}
}

// RunWailsApp runs the Wails application
func RunWailsApp() {
	// Setup file logger
	// Determine log directory - use logs folder relative to executable or backend folder
	logDir := "logs"
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		// If running from build/bin, use logs in build directory
		if filepath.Base(exeDir) == "bin" {
			logDir = filepath.Join(exeDir, "..", "logs")
		} else {
			logDir = filepath.Join(exeDir, "logs")
		}
	}

	// Create file logger with Debug level (will show all logs)
	// Wails logger levels: TRACE=1, DEBUG=2, INFO=3, WARNING=4, ERROR=5
	var fileLogger *logger.FileLogger
	fileLogger, err := logger.NewFileLogger(logDir, 2) // DEBUG level = 2
	if err != nil {
		log.Printf("Failed to create file logger: %v, using console only", err)
		fileLogger = nil
	} else {
		defer fileLogger.Close()
		fileLogger.Info("Starting Wails application...")
		fileLogger.Info("Log directory: " + logDir)
	}

	// Setup standard logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Wails application...")

	// Create an instance of the app structure
	app := NewApp()

	// Store logger in app if needed
	if fileLogger != nil {
		app.fileLogger = fileLogger
	}

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

	// Check if we have embedded assets (production build)
	// Try to check if assets is not empty by checking for index.html
	assetServerOptions := &assetserver.Options{}

	// Try to read from embedded assets first
	_, errEmbed := assets.Open("index.html")
	if errEmbed == nil {
		// Use embedded assets (production build)
		assetServerOptions.Assets = assets
		log.Println("Using embedded assets (production mode)")
	} else {
		log.Printf("Embedded assets not found, error: %v", errEmbed)
		// Try file system path for development
		if _, err := os.Stat(frontendDistPath); err == nil {
			assetServerOptions.Assets = os.DirFS(frontendDistPath)
			log.Printf("Using file system assets from %s (development mode)", frontendDistPath)
		} else {
			// Fallback: try frontend-dist folder in backend
			// For production, try relative to executable
			localFrontendDist := "frontend-dist"

			// Also check relative to executable location
			if exePath, err := os.Executable(); err == nil {
				exeDir := filepath.Dir(exePath)
				// Try multiple locations
				possiblePaths := []string{
					localFrontendDist,
					filepath.Join(exeDir, "frontend-dist"),
					filepath.Join(exeDir, "..", "frontend-dist"),
					filepath.Join(exeDir, "..", "..", "frontend-dist"),
				}

				found := false
				for _, path := range possiblePaths {
					if _, err := os.Stat(path); err == nil {
						assetServerOptions.Assets = os.DirFS(path)
						log.Printf("Using local frontend-dist folder from %s", path)
						found = true
						break
					}
				}

				if !found {
					// Last resort: use embedded assets anyway
					assetServerOptions.Assets = assets
					log.Printf("Warning: Frontend dist not found in any location, using embedded assets")
					log.Printf("Tried paths: %v", possiblePaths)
				}
			} else if _, err := os.Stat(localFrontendDist); err == nil {
				assetServerOptions.Assets = os.DirFS(localFrontendDist)
				log.Printf("Using local frontend-dist folder (development mode)")
			} else {
				// Last resort: use empty embed (will show error)
				assetServerOptions.Assets = assets
				log.Printf("Warning: Frontend dist not found, application may not work correctly")
			}
		}
	}

	// Create application with options
	appOptions := &options.App{
		Title:            "Warehouse Management System",
		Width:            1280,
		Height:           800,
		AssetServer:      assetServerOptions,
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.OnStartup,
		OnDomReady:       app.OnDomReady,
		OnBeforeClose:    app.OnBeforeClose,
		OnShutdown:       app.OnShutdown,
		Logger:           fileLogger, // Use file logger for Wails
		LogLevel:         2,          // DEBUG level = 2
	}

	if fileLogger != nil {
		fileLogger.Info("Running Wails application...")
	}
	log.Println("Running Wails application...")
	err = wails.Run(appOptions)

	if err != nil {
		log.Printf("Fatal error starting Wails application: %v", err)
		// Show error message box on Windows
		// For now, just log to console
		os.Exit(1)
	}
}

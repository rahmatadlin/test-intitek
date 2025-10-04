package main

import (
	"log"
	"warehouse-management/config"
	"warehouse-management/database"
	"warehouse-management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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
			"status": "ok",
			"message": "Warehouse Management API is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}


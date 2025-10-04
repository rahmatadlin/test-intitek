package routes

import (
	"warehouse-management/config"
	"warehouse-management/controllers"
	"warehouse-management/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	authController := &controllers.AuthController{Config: cfg}

	// Public routes
	api := router.Group("/api")
	{
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
		}
	}

	// Protected routes (require authentication)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg))
	{
		// Product routes
		products := protected.Group("/products")
		{
			products.GET("", controllers.GetProducts)
			products.GET("/:id", controllers.GetProduct)
			products.POST("", controllers.CreateProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

		// Dashboard routes
		protected.GET("/dashboard/stats", controllers.GetDashboardStats)

		// Export routes
		protected.GET("/export/csv", controllers.ExportCSV)

		// Barcode routes
		protected.GET("/barcode/:sku", controllers.GenerateBarcode)
	}
}


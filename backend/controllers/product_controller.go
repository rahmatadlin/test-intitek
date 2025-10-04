package controllers

import (
	"bytes"
	"encoding/csv"
	"image/png"
	"net/http"
	"strconv"
	"warehouse-management/database"
	"warehouse-management/models"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/gin-gonic/gin"
)

// GetProducts retrieves all products with optional filtering
func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB

	// Filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by low stock
	if lowStock := c.Query("low_stock"); lowStock == "true" {
		query = query.Where("status = ?", "low_stock")
	}

	if err := query.Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// GetProduct retrieves a single product by ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// CreateProduct adds a new product to the inventory
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Auto-update status based on quantity
	product.UpdateStatus()

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create product. SKU might already exist."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": product})
}

// UpdateProduct modifies an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	product.Name = updateData.Name
	product.SKU = updateData.SKU
	product.Quantity = updateData.Quantity
	product.Location = updateData.Location
	product.Status = updateData.Status

	// Auto-update status based on quantity
	product.UpdateStatus()

	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

// DeleteProduct removes a product from the inventory
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetDashboardStats returns warehouse statistics
func GetDashboardStats(c *gin.Context) {
	var totalProducts int64
	var totalStock int64
	var lowStockCount int64

	database.DB.Model(&models.Product{}).Count(&totalProducts)
	database.DB.Model(&models.Product{}).Select("COALESCE(SUM(quantity), 0)").Scan(&totalStock)
	database.DB.Model(&models.Product{}).Where("status = ?", "low_stock").Count(&lowStockCount)

	var lowStockProducts []models.Product
	database.DB.Where("status = ?", "low_stock").Limit(5).Find(&lowStockProducts)

	c.JSON(http.StatusOK, gin.H{
		"total_products":      totalProducts,
		"total_stock":         totalStock,
		"low_stock_count":     lowStockCount,
		"low_stock_products":  lowStockProducts,
	})
}

// ExportCSV exports all products as a CSV file
func ExportCSV(c *gin.Context) {
	var products []models.Product

	if err := database.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Create CSV buffer
	buffer := new(bytes.Buffer)
	writer := csv.NewWriter(buffer)

	// Write header
	header := []string{"ID", "Name", "SKU", "Quantity", "Location", "Status", "Created At", "Updated At"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV header"})
		return
	}

	// Write data
	for _, product := range products {
		row := []string{
			strconv.Itoa(int(product.ID)),
			product.Name,
			product.SKU,
			strconv.Itoa(product.Quantity),
			product.Location,
			product.Status,
			product.CreatedAt.Format("2006-01-02 15:04:05"),
			product.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV row"})
			return
		}
	}

	writer.Flush()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=products.csv")
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
}

// GenerateBarcode generates a barcode image for a product's SKU
func GenerateBarcode(c *gin.Context) {
	sku := c.Param("sku")

	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SKU is required"})
		return
	}

	// Generate barcode
	barcodeImg, err := code128.Encode(sku)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate barcode"})
		return
	}

	// Scale barcode to 200x100 pixels
	barcodeImg, err = barcode.Scale(barcodeImg, 200, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scale barcode"})
		return
	}

	// Encode to PNG
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, barcodeImg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode barcode"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "inline; filename=barcode-"+sku+".png")
	c.Data(http.StatusOK, "image/png", buffer.Bytes())
}


package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"warehouse-management/config"
	"warehouse-management/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase establishes a connection to the SQLite database
func ConnectDatabase(cfg *config.Config) error {
	// Get database path
	dbPath := cfg.DBPath
	if dbPath == "" {
		// Default to data folder relative to executable or current directory
		dbPath = "data/warehouse.db"
	}

	// Create data directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if dbDir != "." && dbDir != "" {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return fmt.Errorf("failed to create database directory: %v", err)
		}
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	if err := DB.AutoMigrate(&models.Product{}, &models.User{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")

	// Create default admin user if no users exist
	createDefaultUser()

	// Seed sample products
	SeedProducts()

	return nil
}

// createDefaultUser creates a default admin user for testing
func createDefaultUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	log.Printf("Checking for existing users... Found: %d users", count)

	if count == 0 {
		log.Println("No users found, creating default admin user...")
		user := models.User{
			Username: "admin",
			Email:    "admin@warehouse.com",
		}
		if err := user.HashPassword("admin123"); err != nil {
			log.Printf("ERROR: Failed to hash default password: %v", err)
			return
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Printf("ERROR: Failed to create default user: %v", err)
		} else {
			log.Printf("SUCCESS: Default admin user created (username: admin, password: admin123)")
			log.Printf("User ID: %d, Username: %s, Email: %s", user.ID, user.Username, user.Email)
		}
	} else {
		log.Printf("Users already exist (%d), skipping default user creation", count)
		// Log existing users for debugging
		var users []models.User
		if err := DB.Find(&users).Error; err == nil {
			for _, u := range users {
				log.Printf("Existing user: ID=%d, Username=%s, Email=%s", u.ID, u.Username, u.Email)
			}
		}
	}
}

package database

import (
	"fmt"
	"log"
	"warehouse-management/config"
	"warehouse-management/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase establishes a connection to the MySQL database
func ConnectDatabase(cfg *config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	return nil
}

// createDefaultUser creates a default admin user for testing
func createDefaultUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)

	if count == 0 {
		user := models.User{
			Username: "admin",
			Email:    "admin@warehouse.com",
		}
		if err := user.HashPassword("admin123"); err != nil {
			log.Println("Error hashing default password:", err)
			return
		}

		if err := DB.Create(&user).Error; err != nil {
			log.Println("Error creating default user:", err)
		} else {
			log.Println("Default admin user created (username: admin, password: admin123)")
		}
	}
}


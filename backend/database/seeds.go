package database

import (
	"log"
	"warehouse-management/models"
)

// SeedProducts creates sample product data if the products table is empty
func SeedProducts() {
	var count int64
	DB.Model(&models.Product{}).Count(&count)

	// Only seed if there are no products
	if count == 0 {
		products := []models.Product{
			{
				Name:     "Laptop Dell XPS 15",
				SKU:      "LAPTOP-001",
				Quantity: 25,
				Location: "Warehouse A, Shelf 12",
			},
			{
				Name:     "Wireless Mouse Logitech MX Master 3",
				SKU:      "MOUSE-001",
				Quantity: 3,
				Location: "Warehouse A, Shelf 3",
			},
			{
				Name:     "Mechanical Keyboard RGB",
				SKU:      "KEYB-001",
				Quantity: 50,
				Location: "Warehouse B, Shelf 7",
			},
			{
				Name:     "Monitor 27 inch 4K",
				SKU:      "MON-001",
				Quantity: 0,
				Location: "Warehouse A, Shelf 15",
			},
			{
				Name:     "USB-C Hub Multiport",
				SKU:      "USB-001",
				Quantity: 15,
				Location: "Warehouse B, Shelf 2",
			},
		}

		for _, product := range products {
			// Auto-generate status based on quantity
			product.UpdateStatus()
			
			if err := DB.Create(&product).Error; err != nil {
				log.Printf("Error seeding product %s: %v", product.Name, err)
			} else {
				log.Printf("Seeded product: %s (Qty: %d, Status: %s)", product.Name, product.Quantity, product.Status)
			}
		}

		log.Println("Product seeding completed!")
	} else {
		log.Printf("Products already exist (%d), skipping seeding", count)
	}
}


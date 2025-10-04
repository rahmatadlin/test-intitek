package models

import (
	"time"
)

// Product represents a warehouse inventory item
type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null" binding:"required"`
	SKU       string    `json:"sku" gorm:"uniqueIndex;not null" binding:"required"`
	Quantity  int       `json:"quantity" gorm:"not null" binding:"required,min=0"`
	Location  string    `json:"location" gorm:"not null" binding:"required"`
	Status    string    `json:"status" gorm:"not null" binding:"required,oneof=in_stock low_stock out_of_stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name
func (Product) TableName() string {
	return "products"
}

// UpdateStatus automatically updates the status based on quantity
func (p *Product) UpdateStatus() {
	if p.Quantity == 0 {
		p.Status = "out_of_stock"
	} else if p.Quantity <= 10 {
		p.Status = "low_stock"
	} else {
		p.Status = "in_stock"
	}
}


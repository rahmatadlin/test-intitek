package models

import (
	"time"
)

// Product represents a warehouse inventory item
type Product struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null" binding:"required"`
	SKU       string    `json:"sku" gorm:"type:varchar(100);uniqueIndex;not null" binding:"required"`
	Quantity  int       `json:"quantity" gorm:"not null" binding:"required,min=0"`
	Location  string    `json:"location" gorm:"type:varchar(255);not null" binding:"required"`
	Status    string    `json:"status" gorm:"type:varchar(50);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name
func (Product) TableName() string {
	return "products"
}

// UpdateStatus automatically updates the status based on quantity
// Rules: 0 = out_of_stock, 1-5 = low_stock, >5 = in_stock
func (p *Product) UpdateStatus() {
	if p.Quantity == 0 {
		p.Status = "out_of_stock"
	} else if p.Quantity <= 5 {
		p.Status = "low_stock"
	} else {
		p.Status = "in_stock"
	}
}


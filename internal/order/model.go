package order

import (
	"ecommerce/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   int               `gorm:"column:user_id;not null;index:user_index"`
	Code     string            `gorm:"size:100;unique;not null"`
	Price    float64           `gorm:"not null;"`
	Status   string            `gorm:"size:50;not null"`
	Products []product.Product `gorm:"many2many:orders_products;"`
}

type OrderWithProducts struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id" binding:"required"`
	ProductsID []int   `json:"products_id" binding:"required"`
	Code       string  `json:"code"    binding:"required"`
	Price      float64 `json:"price"   binding:"required"`
	Status     string  `json:"status"  binding:"required"`
}

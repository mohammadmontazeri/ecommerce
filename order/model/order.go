package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID int     `gorm:"column:user_id;not null;index:user_index"`
	Code   string  `gorm:"size:100;unique;not null"`
	Price  float64 `gorm:"not null;"`
	Status string  `gorm:"size:50;not null"`
}

type OrderWithProducts struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id" binding:"required"`
	ProductsID []int   `json:"products_id" binding:"required"`
	Code       string  `json:"code"    binding:"required"`
	Price      float64 `json:"price"   binding:"required"`
	Status     string  `json:"status"  binding:"required"`
}

type OrderInput struct {
	Order      Order
	ProductsID []int
}

type OrderService interface {
	Create(input OrderWithProducts) error
	Read(orderID int) (OrderWithProducts, error)
	Update(orderInput OrderWithProducts, orderID int) error
	Delete(orderID int) error
}

type OrderRepository interface {
	InsertOrderWithoutProducts(order Order) (Order, error)
	AddOrderToPivotTable(productsID []int, order Order) error
	GetOrderFromId(intQueryParameter int) (Order, error)
	GetOrderProducts(orderWithoutProduct Order) ([]int, error)
	UpdateOrderRow(input Order, orderID int) (Order, error)
	DeleteOrderProduct(orderID int) error
	DeleteRow(order Order, orderID int) error
}

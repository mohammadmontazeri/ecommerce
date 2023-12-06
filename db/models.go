package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `gorm:"column:username;unique;" json:"username"`
	Email    string  `gorm:"column:email;unique;not null;index:em_index" json:"email"`
	Password string  `gorm:"column:password;not null" json:"password"`
	Orders   []Order `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Category struct {
	gorm.Model
	Name     string `gorm:"not null;unique;size:255"`
	ParentID *int   `gorm:"not null" json:"parent_id"`
	Parent   *Category
	Products []Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Product struct {
	gorm.Model
	Code       string  `gorm:"size:100;unique;not null" form:"code"`
	Title      string  `gorm:"size:255;unique;not null" form:"title"`
	Price      float64 `gorm:"not null" form:"price"`
	Picture    string
	Detail     string `form:"detail"`
	CategoryID uint   `gorm:"column:category_id;not null;index:catgeory_index" form:"category_id"`
}

type Order struct {
	gorm.Model
	UserID   int       `gorm:"column:user_id;not null;index:user_index"`
	Code     string    `gorm:"size:100;unique;not null"`
	Price    float64   `gorm:"not null;"`
	Status   string    `gorm:"size:50;not null"`
	Products []Product `gorm:"many2many:orders_products;"`
}

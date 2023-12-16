package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Code       string  `gorm:"size:100;unique;not null" form:"code"`
	Title      string  `gorm:"size:255;unique;not null" form:"title"`
	Price      float64 `gorm:"not null" form:"price"`
	Picture    string
	Detail     string `form:"detail"`
	CategoryID uint   `gorm:"column:category_id;not null;index:catgeory_index" form:"category_id"`
}

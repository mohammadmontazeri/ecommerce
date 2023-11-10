package product

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code       string  `gorm:"size:100;unique;not null" form:"code"`
	Title      string  `gorm:"size:255;unique;not null" form:"title"`
	Price      float64 `gorm:"not null" form:"price"`
	Picture    string
	Detail     string `form:"detail"`
	CategoryID uint   `gorm:"column:category_id;not null;index:catgeory_index" form:"category_id"`
}

type ProductInput struct {
	Id         int                   `json:"id"`
	Code       string                `form:"code" json:"code" binding:"required"`
	Title      string                `form:"title" json:"title" binding:"required"`
	Price      float64               `form:"price" json:"price" binding:"required"`
	Picture    *multipart.FileHeader `form:"picture" binding:"required"`
	Detail     string                `form:"detail" json:"detail" binding:"required"`
	CategoryID uint                  `form:"category_id" json:"category_id" binding:"required"`
}

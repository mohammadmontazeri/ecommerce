package category

import (
	"ecommerce/internal/product"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string `gorm:"not null;unique;size:255"`
	ParentID *int   `gorm:"not null" json:"parent_id"`
	Parent   *Category
	Products []product.Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

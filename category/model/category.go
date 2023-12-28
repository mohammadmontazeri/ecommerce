package model

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name     string `gorm:"not null;unique;size:255"`
	ParentID *int   `gorm:"not null" json:"parent_id"`
	Parent   *Category
}

type CategoryService interface {
	Create(input Category) error
	Read(categoryID int) (Category, error)
	Update(categoryID int, input Category) error
	Delete(categoryID int) error
}

type CategoryRepository interface {
	InsertCategory(input Category) error
	GetCategory(category Category, categoryID int) (Category, error)
	UpdateRow(category Category, categoryID int) error
	DeleteRow(category Category, categoryID int) error
}

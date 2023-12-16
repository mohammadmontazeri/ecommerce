package category

import (
	"ecommerce/models"
	"errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) models.CategoryRepository {
	return &categoryRepository{
		DB: db,
	}
}

func (o *categoryRepository) InsertCategory(category models.Category) error {

	res := o.DB.Create(&category)

	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}

func (o *categoryRepository) GetCategory(category models.Category, categoryID int) (models.Category, error) {

	res := o.DB.Find(&category, categoryID)
	if res.Error != nil {
		return category, res.Error
	}

	return category, nil
}

func (o *categoryRepository) UpdateRow(category models.Category, categoryID int) error {
	res := o.DB.Model(&category).Where("id", categoryID).Updates(category)
	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return res.Error
	}

	return nil
}

func (o *categoryRepository) DeleteRow(category models.Category, categoryID int) error {

	res := o.DB.Delete(&category, categoryID)
	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return errors.New("error:Delete not completed ! ")
	} else {
		return nil
	}

}

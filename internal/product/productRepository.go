package product

import (
	"ecommerce/models"
	"errors"

	"gorm.io/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) models.ProductRepository {
	return &productRepository{
		DB: db,
	}
}

func (o *productRepository) InsertProduct(product models.Product) error {

	res := o.DB.Create(&product)

	if res.Error != nil {
		return res.Error
	} else {
		return nil
	}
}

func (o *productRepository) UpdateRow(product models.Product, productID int) error {

	res := o.DB.Model(&product).Where("id", productID).Updates(product)

	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return errors.New("error:No Row found !")
	}
	return nil
}

func (o *productRepository) DeleteRow(product models.Product, productID int) error {

	res := o.DB.Delete(&product, productID)

	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return errors.New("error:No Row found !")
	}

	return nil
}

func (o *productRepository) GetProduct(product models.Product, productID int) (models.Product, error) {
	res := o.DB.Find(&product, productID)

	if res.Error != nil {
		return product, res.Error
	} else {
		return product, nil
	}
}

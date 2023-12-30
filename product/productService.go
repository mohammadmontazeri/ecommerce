package product

import (
	"ecommerce/product/model"
)

type productService struct {
	productService model.ProductRepository
}

func NewProductService(s model.ProductRepository) *productService {
	return &productService{
		productService: s,
	}
}

func (ps *productService) Create(input model.ProductInput, picturePath string) error {

	product := model.Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.CategoryID = input.CategoryID
	product.Picture = picturePath
	err := ps.productService.InsertProduct(product)

	if err != nil {
		return err
	} else {
		return nil
	}

}

func (ps *productService) Update(productID int, input model.ProductInput, picturePath string) error {

	product := model.Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.CategoryID = input.CategoryID
	product.Picture = picturePath

	err := ps.productService.UpdateRow(product, productID)

	if err != nil {
		return err
	} else {
		return nil
	}

}

func (ps *productService) Delete(productID int) error {

	product := model.Product{}
	err := ps.productService.DeleteRow(product, productID)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (ps *productService) Read(productID int) (model.Product, error) {
	product := model.Product{}
	product, err := ps.productService.GetProduct(product, productID)

	if err != nil {
		return product, err
	} else {
		return product, nil
	}
}

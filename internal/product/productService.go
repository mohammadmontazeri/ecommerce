package product

import "ecommerce/models"

type productService struct {
	productService models.ProductRepository
}

func NewProductService(s models.ProductRepository) *productService {
	return &productService{
		productService: s,
	}
}

func (ps *productService) Create(input models.ProductInput, picturePath string) error {

	product := models.Product{}

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

func (ps *productService) Update(productID int, input models.ProductInput, picturePath string) error {

	product := models.Product{}

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

	product := models.Product{}
	err := ps.productService.DeleteRow(product, productID)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (ps *productService) Read(productID int) (models.Product, error) {
	product := models.Product{}
	product, err := ps.productService.GetProduct(product, productID)

	if err != nil {
		return product, err
	} else {
		return product, nil
	}
}

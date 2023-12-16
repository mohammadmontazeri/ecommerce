package order

import (
	"ecommerce/models"
	"errors"
)

type orderService struct {
	orderRepository models.OrderRepository
}

func NewOrderService(s models.OrderRepository) *orderService {
	return &orderService{
		orderRepository: s,
	}
}

func (os *orderService) Create(input models.OrderWithProducts) error {

	if len(input.ProductsID) == 0 {
		return errors.New("error :ProductsID not Corrected !")
	}
	o := models.OrderInput{}

	o.Order.Code = input.Code
	o.Order.UserID = input.UserID
	o.Order.Price = input.Price
	o.Order.Status = input.Status
	o.ProductsID = input.ProductsID

	orderWithoutProduct, err := os.orderRepository.InsertOrderWithoutProducts(o.Order)
	if err != nil {
		return err
	}

	err = os.orderRepository.AddOrderToPivotTable(o.ProductsID, orderWithoutProduct)
	if err != nil {
		return err
	}

	return nil
}

func (os *orderService) Read(orderID int) (models.OrderWithProducts, error) {
	var orderWithoutProduct models.Order
	var orderWithProducts models.OrderWithProducts
	orderWithoutProduct, err := os.orderRepository.GetOrderFromId(orderID)
	if err != nil {
		return orderWithProducts, err
	}

	var orderProducts []int
	orderProducts, err = os.orderRepository.GetOrderProducts(orderWithoutProduct)

	if err != nil {
		return orderWithProducts, err
	}

	order := models.OrderWithProducts{}
	order.ID = int(orderWithoutProduct.ID)
	order.ProductsID = orderProducts
	order.Code = orderWithoutProduct.Code
	order.Status = orderWithoutProduct.Status
	order.Price = orderWithoutProduct.Price
	order.UserID = orderWithoutProduct.UserID

	return order, nil

}

func (os *orderService) Update(input models.OrderWithProducts, orderID int) error {

	order := models.Order{}
	if len(input.ProductsID) == 0 {
		return errors.New("error :ProductsID not Corrected !")
	}

	order.ID = uint(orderID)
	order.Code = input.Code
	order.UserID = input.UserID
	order.Price = input.Price
	order.Status = input.Status

	order, err := os.orderRepository.UpdateOrderRow(order, orderID)

	if err != nil {
		return err
	}

	err = os.orderRepository.DeleteOrderProduct(orderID)

	if err != nil {
		return err
	}

	err = os.orderRepository.AddOrderToPivotTable(input.ProductsID, order)
	if err != nil {
		return nil
	}

	return nil
}

func (os *orderService) Delete(orderID int) error {

	err := os.orderRepository.DeleteOrderProduct(orderID)

	if err != nil {
		return err
	}

	order := models.Order{}
	err = os.orderRepository.DeleteRow(order, orderID)

	if err != nil {
		return err
	}

	return nil
}

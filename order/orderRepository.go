package order

import (
	"ecommerce/order/model"
	"errors"

	"gorm.io/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) model.OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (o *orderRepository) InsertOrderWithoutProducts(order model.Order) (model.Order, error) {

	if res := o.DB.Create(&order); res.Error != nil {
		return order, res.Error
	}
	return order, nil
}

func (o *orderRepository) AddOrderToPivotTable(productsID []int, order model.Order) error {
	for _, productID := range productsID {
		res := o.DB.Exec("INSERT INTO orders_products (ORDER_ID,PRODUCT_ID) VALUES ($1,$2)", order.ID, productID)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func (o *orderRepository) GetOrderFromId(intQueryParameter int) (model.Order, error) {
	var order model.Order
	res := o.DB.Find(&order, intQueryParameter)

	if res.Error != nil {
		return order, res.Error
	}
	return order, nil
}

func (o *orderRepository) GetOrderProducts(orderWithoutProduct model.Order) ([]int, error) {
	var products []int
	rows, err := o.DB.Raw("SELECT product_id FROM orders_products WHERE order_id=$1", orderWithoutProduct.ID).Rows()
	if err != nil {
		return products, err
	}
	for rows.Next() {
		var productID int
		err := rows.Scan(&productID)
		if err != nil {
			return products, err
		}
		products = append(products, productID)
	}
	return products, nil
}

func (o *orderRepository) UpdateOrderRow(order model.Order, orderID int) (model.Order, error) {

	res := o.DB.Model(&order).Where("id", orderID).Updates(order)

	if res.Error != nil {
		return order, res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return order, errors.New("error : Update not completed !")
	}

	return order, nil
}

func (o *orderRepository) DeleteOrderProduct(orderID int) error {

	res := o.DB.Exec("DELETE FROM orders_products WHERE order_id=$1 ;", orderID)
	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected

	if rowsAffected == 0 {
		return errors.New("error :no row found")
	}
	return nil
}

func (o *orderRepository) DeleteRow(order model.Order, orderID int) error {

	res := o.DB.Delete(&order, orderID)

	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		return errors.New("error :Delete not completed !")
	}
	return nil
}

package order

import (
	"ecommerce/db"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderWithProducts struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id" binding:"required"`
	ProductsID []int   `json:"products_id" binding:"required"`
	Code       string  `json:"code"    binding:"required"`
	Price      float64 `json:"price"   binding:"required"`
	Status     string  `json:"status"  binding:"required"`
}

type Order struct {
	db.Order
}

func NewOrderModel(db *gorm.DB) OrderModel {
	return OrderModel{db: db}
}

type OrderModel struct {
	db *gorm.DB
}

func (om *OrderModel) createOrder(c *gin.Context) {
	var tx = om.db.Begin()
	var input OrderWithProducts
	defer tx.Rollback()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.ProductsID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "products are required"})
		return
	}
	order := Order{}

	order.Code = input.Code
	order.UserID = input.UserID
	order.Price = input.Price
	order.Status = input.Status

	order, err := insertOrderWithoutProducts(tx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = addOrderToPivotTable(tx, input.ProductsID, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	commit := tx.Commit()

	if commit.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": commit.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "order add successful"})
	}

}

func (om *OrderModel) Read(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderWithoutProduct Order
	orderWithoutProduct, err = getOrderFromId(om, intQueryParameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var orderProducts []int
	orderProducts, err = getOrderProducts(om, orderWithoutProduct)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query parmeter not set"})
		return
	}

	order := OrderWithProducts{}
	order.ID = int(orderWithoutProduct.ID)
	order.ProductsID = orderProducts
	order.Code = orderWithoutProduct.Code
	order.Status = orderWithoutProduct.Status
	order.Price = orderWithoutProduct.Price
	order.UserID = orderWithoutProduct.UserID

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"order": order})
	}

}

func (om *OrderModel) Update(c *gin.Context) {
	var tx = om.db.Begin() // for db trnsaction
	defer tx.Rollback()
	var input OrderWithProducts

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.ProductsID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "products are required field"})
		return
	}

	order := Order{}

	order.Code = input.Code
	order.UserID = input.UserID
	order.Price = input.Price
	order.Status = input.Status

	res := om.db.Model(&order).Where("id", id).Updates(order)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	err = deleteOrderProduct(om, id)
	if err != nil {
		if err.Error() == "no row found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	err = addOrderToPivotTable(tx, input.ProductsID, order)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}

}

func (om *OrderModel) Delete(c *gin.Context) {
	var tx = om.db.Begin() // for db trnsaction
	defer tx.Rollback()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = deleteOrderProduct(om, id)
	if err != nil {
		if err.Error() == "no row found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// delete order
	res := om.db.Delete(&Order{}, id)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})

}

func deleteOrderProduct(om *OrderModel, id int) error {

	res := om.db.Exec("DELETE FROM orders_products WHERE order_id=$1 ;", id)
	if res.Error != nil {
		return res.Error
	}
	rowsAffected := res.RowsAffected

	if rowsAffected == 0 {
		return errors.New("no row found")
	}
	return nil
}

func addOrderToPivotTable(tx *gorm.DB, productsID []int, order Order) error {

	for _, productID := range productsID {
		res := tx.Exec("INSERT INTO orders_products (ORDER_ID,PRODUCT_ID) VALUES ($1,$2)", order.ID, productID)
		if res.Error != nil {
			return res.Error
		}
	}

	return nil
}

func getOrderFromId(om *OrderModel, id int) (Order, error) {
	var order Order
	res := om.db.Find(&order, id)

	if res.Error != nil {
		return order, res.Error
	}
	return order, nil
}

func getOrderProducts(om *OrderModel, order Order) ([]int, error) {
	var products []int
	rows, err := om.db.Raw("SELECT product_id FROM orders_products WHERE order_id=$1", order.ID).Rows()
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

func insertOrderWithoutProducts(tx *gorm.DB, order Order) (Order, error) {

	if res := tx.Create(&order); res.Error != nil {
		return order, res.Error
	}
	return order, nil
}

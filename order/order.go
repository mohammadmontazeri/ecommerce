package order

import (
	"context"
	"ecommerce/db"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Order struct {
	Id          int     `json:"id"`
	User_id     int     `json:"user_id" binding:"required"`
	Products_id []int   `json:"products_id" binding:"required"`
	Code        string  `json:"code"    binding:"required"`
	Price       float64 `json:"price"   binding:"required"`
	Status      string  `json:"status"  binding:"required"`
}

type GetOrder struct {
	Id      int     `json:"id"`
	User_id int     `json:"user_id" binding:"required"`
	Code    string  `json:"code"    binding:"required"`
	Price   float64 `json:"price"   binding:"required"`
	Status  string  `json:"status"  binding:"required"`
}
type OrderProducts struct {
	Product_id int
}

type GetOrderWithProduct struct {
	Order         GetOrder
	OrderProducts []OrderProducts
}

var DB = db.ConnectToDb()
var ctx = context.Background()
var tx, _ = DB.BeginTx(ctx, nil) // for db trnsaction

func Create(c *gin.Context) {
	var input Order

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order := Order{}

	order.Code = input.Code
	order.User_id = input.User_id
	order.Price = input.Price
	order.Status = input.Status
	_, err := tx.ExecContext(ctx, `INSERT INTO orders(code,user_id,price,status) VALUES($1,$2,$3,$4)`, order.Code, order.User_id, order.Price, order.Status)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	orderValue, err := getOrderFromCode(order.Code, c)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	addOrderToPivotTable(input.Products_id, orderValue.Id, c)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	tx.Commit()
}

func Read(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, _ := strconv.Atoi(queryParameter)

	if queryParameter == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "query parmeter not set"})
	}

	// get order without products
	var orderWithoutProduct GetOrder
	orderWithoutProduct, _ = getOrderFromId(intQueryParameter, c)

	// get order products
	var orderProducts []OrderProducts
	orderProducts, _ = getOrderProducts(intQueryParameter, c)

	var order GetOrderWithProduct

	order = GetOrderWithProduct{
		Order:         orderWithoutProduct,
		OrderProducts: orderProducts,
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"ids": order})

}

func addOrderToPivotTable(products_id []int, order_id int, c *gin.Context) error {
	var err error
	for _, value := range products_id {
		_, err = tx.ExecContext(ctx, "INSERT INTO orders_products(product_id,order_id) VALUES($1,$2)", value, order_id)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "order add successful"})

	return nil
}

func getOrderFromCode(code string, c *gin.Context) (GetOrder, error) {
	var order GetOrder
	err := tx.QueryRowContext(ctx, "SELECT id,user_id,code,price,status FROM orders WHERE code=$1 ", code).Scan(&order.Id, &order.User_id, &order.Code, &order.Price, &order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return order, err
	}
	return order, nil
}

func getOrderFromId(id int, c *gin.Context) (GetOrder, error) {
	var order GetOrder
	err := tx.QueryRowContext(ctx, "SELECT id,user_id,code,price,status FROM orders WHERE id=$1 ", id).Scan(&order.Id, &order.User_id, &order.Code, &order.Price, &order.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return order, err
	}
	return order, nil
}

func getOrderProducts(order_id int, c *gin.Context) ([]OrderProducts, error) {
	products := []OrderProducts{}
	rows, err := tx.QueryContext(ctx, "SELECT product_id FROM orders_products WHERE order_id=$1", order_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return products, err
	}
	for rows.Next() {
		o := OrderProducts{}
		err = rows.Scan(&o.Product_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return products, err
		}
		products = append(products, o)
	}
	return products, nil
}

package order

import (
	"context"
	"database/sql"
	"ecommerce/db"
	"errors"
	"fmt"
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

func Create(c *gin.Context) {
	var tx, err = DB.BeginTx(ctx, nil)
	var input Order
	defer tx.Rollback()

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.Products_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "products are required"})
		return
	}
	order := Order{}

	order.Code = input.Code
	order.User_id = input.User_id
	order.Price = input.Price
	order.Status = input.Status

	// insert order without products
	err = insertOrderWithoutProducts(tx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// get order for its id
	orderValue, err := getOrderFromCode(tx, order.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// insert realtions in pivot table
	err = addOrderToPivotTable(tx, input.Products_id, orderValue.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "order add successful"})
	}

}

func Read(c *gin.Context) {
	var tx, _ = DB.BeginTx(ctx, nil) // for db trnsaction
	defer tx.Rollback()
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get order without products
	var orderWithoutProduct GetOrder
	orderWithoutProduct, err = getOrderFromId(tx, intQueryParameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// get order products
	var orderProducts []OrderProducts
	orderProducts, err = getOrderProducts(tx, intQueryParameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query parmeter not set"})
		return
	}

	var order GetOrderWithProduct

	order = GetOrderWithProduct{
		Order:         orderWithoutProduct,
		OrderProducts: orderProducts,
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"order": order})
	}

}

func Update(c *gin.Context) {
	var tx, _ = DB.BeginTx(ctx, nil) // for db trnsaction
	defer tx.Rollback()
	var input Order

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input.Products_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "products are required field"})
		return
	}

	order := Order{}

	order.Code = input.Code
	order.User_id = input.User_id
	order.Price = input.Price
	order.Status = input.Status

	res, err := tx.ExecContext(ctx, "UPDATE orders SET code=$1,user_id=$2,price=$3,status=$4 WHERE id=$5", order.Code, order.User_id, order.Price, order.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	// delete expired order products
	err = deleteOrderProduct(tx, id)
	if err != nil {
		if err.Error() == "no row found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	// insert order products to pivot table after delete
	err = addOrderToPivotTable(tx, input.Products_id, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}

}

func Delete(c *gin.Context) {
	var tx, _ = DB.BeginTx(ctx, nil) // for db trnsaction
	defer tx.Rollback()
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// delete order products
	err = deleteOrderProduct(tx, id)
	if err != nil {
		if err.Error() == "no row found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// delete order
	res, err := tx.ExecContext(ctx, "DELETE FROM orders WHERE id=$1 ;", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})
	}

}

func deleteOrderProduct(tx *sql.Tx, id int) error {
	res, err := tx.ExecContext(ctx, "DELETE FROM orders_products WHERE order_id=$1 ;", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no row found")
	}
	return nil
}

func addOrderToPivotTable(tx *sql.Tx, products_id []int, order_id int) error {
	fmt.Println(products_id)
	var err error
	for _, value := range products_id {
		if _, err = tx.ExecContext(ctx, "INSERT INTO orders_products(product_id,order_id) VALUES($1,$2)", value, order_id); err != nil {
			return err
		}
	}

	return nil
}

func getOrderFromCode(tx *sql.Tx, code string) (GetOrder, error) {
	var order GetOrder
	if err := tx.QueryRowContext(ctx, "SELECT id,user_id,code,price,status FROM orders WHERE code=$1 ", code).Scan(&order.Id, &order.User_id, &order.Code, &order.Price, &order.Status); err != nil {
		return order, err
	}
	return order, nil
}

func getOrderFromId(tx *sql.Tx, id int) (GetOrder, error) {
	var order GetOrder
	err := tx.QueryRowContext(ctx, "SELECT id,user_id,code,price,status FROM orders WHERE id=$1 ", id).Scan(&order.Id, &order.User_id, &order.Code, &order.Price, &order.Status)
	if err != nil {
		return order, err
	}
	return order, nil
}

func getOrderProducts(tx *sql.Tx, order_id int) ([]OrderProducts, error) {
	products := []OrderProducts{}
	rows, err := tx.QueryContext(ctx, "SELECT product_id FROM orders_products WHERE order_id=$1", order_id)

	if err != nil {
		return products, err
	}
	for rows.Next() {
		o := OrderProducts{}
		err = rows.Scan(&o.Product_id)
		if err != nil {
			return products, err
		}
		products = append(products, o)
	}
	return products, nil
}

func insertOrderWithoutProducts(tx *sql.Tx, order Order) error {

	if _, err := tx.ExecContext(ctx, `INSERT INTO orders(code,user_id,price,status) VALUES($1,$2,$3,$4)`, order.Code, order.User_id, order.Price, order.Status); err != nil {
		return err
	}
	return nil
}

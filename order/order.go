package order

import (
	// "context"
	"database/sql"
	"ecommerce/db"
	"errors"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderWithProducts struct {
	Id         int     `json:"id"`
	UserID     int     `json:"user_id" binding:"required"`
	ProductsID []int   `json:"products_id" binding:"required"`
	Code       string  `json:"code"    binding:"required"`
	Price      float64 `json:"price"   binding:"required"`
	Status     string  `json:"status"  binding:"required"`
}

type Order struct {
	Id     int     `json:"id"`
	UserID int     `json:"user_id" binding:"required"`
	Code   string  `json:"code"    binding:"required"`
	Price  float64 `json:"price"   binding:"required"`
	Status string  `json:"status"  binding:"required"`
}
type OrderProducts struct {
	ProductID int
}

type GetOrderWithProduct struct {
	Order         Order
	OrderProducts []OrderProducts
}

type Connector interface {
	ConnectDB() *gorm.DB
}

func (cm OrderModel) ConnectDB() *gorm.DB {
	return db.ConnectToDBGorm()
}
func NewStruct(c Connector) *OrderModel {
	return &OrderModel{connector: c}
}

type OrderModel struct {
	connector Connector
}

func (om *OrderModel) Create(c *gin.Context) {
	var tx = om.connector.ConnectDB().Begin()
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
	order := db.Order{}

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

func Read(c *gin.Context) {
	var tx, _ = DB.BeginTx(ctx, nil) // for db trnsaction
	defer tx.Rollback()
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var orderWithoutProduct Order
	orderWithoutProduct, err = getOrderFromId(tx, intQueryParameter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	res, err := tx.ExecContext(ctx, "UPDATE orders SET code=$1,user_id=$2,price=$3,status=$4 WHERE id=$5", order.Code, order.UserID, order.Price, order.Status, id)
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

	err = deleteOrderProduct(tx, id)
	if err != nil {
		if err.Error() == "no row found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	err = addOrderToPivotTable(tx, input.ProductsID, id)

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

func addOrderToPivotTable(tx *gorm.DB, productsID []int, order db.Order) error {

	err := tx.Model(&order).Association("Products").Append([]db.Product{})

	if err != nil {
		return err
	}

	return nil
}

func getOrderFromId(tx *sql.Tx, id int) (Order, error) {
	var order Order
	err := tx.QueryRowContext(ctx, "SELECT id,user_id,code,price,status FROM orders WHERE id=$1 ", id).Scan(&order.Id, &order.UserID, &order.Code, &order.Price, &order.Status)
	if err != nil {
		return order, err
	}
	return order, nil
}

func getOrderProducts(tx *sql.Tx, orderID int) ([]OrderProducts, error) {
	products := []OrderProducts{}
	rows, err := tx.QueryContext(ctx, "SELECT product_id FROM orders_products WHERE order_id=$1", orderID)

	if err != nil {
		return products, err
	}
	for rows.Next() {
		o := OrderProducts{}
		err = rows.Scan(&o.ProductID)
		if err != nil {
			return products, err
		}
		products = append(products, o)
	}
	return products, nil
}

func insertOrderWithoutProducts(tx *gorm.DB, order db.Order) (db.Order, error) {

	if res := tx.Create(&order); res.Error != nil {
		return order, res.Error
	}
	return order, nil
}

package main

import (
	"ecommerce/auth"
	"ecommerce/db"
	"ecommerce/internal/category"

	"ecommerce/internal/order"
	"ecommerce/internal/product"
	"ecommerce/internal/user"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()
	var userWithDB = user.New(DB)
	public := r.Group("/api")
	public.POST("/register", userWithDB.Register)
	public.POST("/login", userWithDB.Login)

	protected := r.Group("/api/admin")
	protected.Use(auth.JwtApiMiddleware)
	protected.GET("/user", user.AuthorizedUser)
	// category crud
	var categoryWithDB = category.New(DB)
	protected.POST("/category/create", categoryWithDB.Create)
	protected.GET("/category/:id", categoryWithDB.Read)
	protected.PUT("/category/update/:id", categoryWithDB.Update)
	protected.DELETE("/category/delete/:id", categoryWithDB.Delete)
	// product crud
	var prodcutWithDB = product.New(DB)

	protected.POST("/product/create", prodcutWithDB.Create)
	protected.GET("/product/:id", prodcutWithDB.Read)
	protected.PUT("/product/update/:id", prodcutWithDB.Update)
	protected.DELETE("/product/delete/:id", prodcutWithDB.Delete)
	// order crud
	var orderModel = order.NewModel(DB)
	var orderWithDB = order.NewOrderModel(&orderModel)
	protected.POST("/order/create", orderWithDB.CreateOrder)
	// protected.GET("/order/:id", orderWithDB.Read)
	// protected.PUT("/order/update/:id", orderWithDB.Update)
	// protected.DELETE("/order/delete/:id", orderWithDB.Delete)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

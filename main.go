package main

import (
	"ecommerce/auth"
	"ecommerce/category"
	"ecommerce/db"
	"ecommerce/order"
	"ecommerce/product"
	"ecommerce/user"

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
	var orderWithDB = order.New(DB)
	protected.POST("/order/create", orderWithDB.Create)
	protected.GET("/order/:id", orderWithDB.Read)
	protected.PUT("/order/update/:id", orderWithDB.Update)
	protected.DELETE("/order/delete/:id", orderWithDB.Delete)

	r.Run(":8080")

}

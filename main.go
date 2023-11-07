package main

import (
	"ecommerce/auth"
	"ecommerce/category"
	"ecommerce/db"
	"ecommerce/order"
	"ecommerce/product"
	"ecommerce/user"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {


	r := gin.Default()

	public := r.Group("/api")
	public.GET("/migrate", db.MigrateTables)
	public.POST("/register", user.NewStruct(user.UserModel{}).Register)
	public.POST("/login", user.NewStruct(user.UserModel{}).Login)

	protected := r.Group("/api/admin")
	protected.Use(auth.JwtApiMiddleware)
	protected.GET("/user", user.AuthorizedUser)
	// category crud
	protected.POST("/category/create", category.NewStruct(category.CategoryModel{}).Create)
	protected.GET("/category/:id", category.NewStruct(category.CategoryModel{}).Read)
	protected.PUT("/category/update/:id", category.NewStruct(category.CategoryModel{}).Update)
	protected.DELETE("/category/delete/:id", category.NewStruct(category.CategoryModel{}).Delete)
	// product crud
	protected.POST("/product/create", product.Create)
	protected.GET("/product/:id", product.Read)
	protected.PUT("/product/update/:id", product.Update)
	protected.DELETE("/product/delete/:id", product.Delete)
	// order crud
	protected.POST("/order/create", order.NewStruct(order.OrderModel{}).Create)
	protected.GET("/order/:id", order.Read)
	protected.PUT("/order/update/:id", order.Update)
	protected.DELETE("/order/delete/:id", order.Delete)

	r.Run(":8080")

}

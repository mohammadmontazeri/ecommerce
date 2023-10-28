package main

import (
	"ecommerce/auth"
	"ecommerce/category"
	"ecommerce/db"
	"ecommerce/product"
	"ecommerce/user"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	r := gin.Default()

	public := r.Group("/api")
	public.GET("/migrate", db.MigrateTables)
	public.POST("/register", user.Register)
	public.POST("/login", user.Login)

	protected := r.Group("/api/admin")
	protected.Use(auth.JwtApiMiddleware)
	protected.GET("/user", user.AuthorizedUser)
	// category crud
	protected.POST("/category/create", category.Create)
	protected.GET("/category/:id", category.Read)
	protected.PUT("/category/update/:id", category.Update)
	protected.DELETE("/category/delete/:id", category.Delete)
	// product crud
	protected.POST("/product/create", product.Create)
	protected.GET("/product/:id", product.Read)
	protected.PUT("/product/update/:id", product.Update)
	protected.DELETE("/product/delete/:id", product.Delete)

	r.Run(":8080")

}

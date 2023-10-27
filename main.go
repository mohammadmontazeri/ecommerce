package main

import (
	"ecommerce/auth"
	"ecommerce/category"
	"ecommerce/db"
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
	public.POST("/category/create", category.Create)
	public.GET("/category/:id", category.Read)
	protected.PUT("/category/update/:id", category.Update)
	public.DELETE("/category/delete/:id", category.Delete)
	// // product crud
	// public.POST("/product/create", api.CreateProduct)
	// public.GET("/product/:id", api.GetProduct)
	// public.PUT("/product/update/:id", api.UpdateProduct)
	// public.DELETE("/product/delete/:id", api.DeleteProduct)

	// public.POST("/upload",api.Upload)
	r.Run(":8080")

}

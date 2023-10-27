package main

import (
	"ecommerce/auth"
	"ecommerce/db"
	"ecommerce/user"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	db.ConnectToDb()
	r := gin.Default()
	public := r.Group("/api")
	public.GET("/migrate", db.MigrateTables)
	public.POST("/register", user.Register)
	public.POST("/login", user.Login)

	protected := r.Group("/api/admin")
	protected.Use(auth.JwtApiMiddleware())
	protected.GET("/user", user.AuthorizedUser)
	// category crud
	// public.POST("/category/create", api.CreateCategory)
	// public.GET("/category/:id", api.GetCategory)
	// public.PUT("/category/update/:id", api.UpdateCategory)
	// public.DELETE("/category/delete/:id", api.DeleteCategory)
	// // product crud
	// public.POST("/product/create", api.CreateProduct)
	// public.GET("/product/:id", api.GetProduct)
	// public.PUT("/product/update/:id", api.UpdateProduct)
	// public.DELETE("/product/delete/:id", api.DeleteProduct)

	// public.POST("/upload",api.Upload)
	r.Run(":8080")

}

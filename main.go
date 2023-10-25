package main

import (
	// "database/sql"

	// "fmt"
	"main/api"
	"main/middleware"
	"main/model"

	// "main/model"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	// tables := map[string]string{
	// 	"userTable" : "userTdssdsable" ,
	// 	"categoryTable" : "ddsds" ,
	// 	"productTable" : "qweqwe" ,
	// }  

	// for _,value := range tables {
	// 	// db.Exec(value)
	// 	fmt.Println(value)
	// } 
	model.ConnectToDb()
	r := gin.Default()
	public := r.Group("/api")
	public.GET("/migrate",model.MigrateTables)
	public.POST("/register", api.Register)
	public.POST("/login",api.Login)

	protected := r.Group("/api/admin")
	protected.Use(middleware.JwtApiMiddleware())
	protected.GET("/user",api.AuthorizedUser)
	// category crud
	public.POST("/category/create", api.CreateCategory)
	public.GET("/category/:id",api.GetCategory)
	public.PUT("/category/update/:id",api.UpdateCategory)
	public.DELETE("/category/delete/:id",api.DeleteCategory)
	// product crud
	public.POST("/product/create", api.CreateProduct)
	public.GET("/product/:id",api.GetProduct)
	public.PUT("/product/update/:id",api.UpdateProduct)
	public.DELETE("/product/delete/:id",api.DeleteProduct)

	// public.POST("/upload",api.Upload)
	r.Run(":8080")

}

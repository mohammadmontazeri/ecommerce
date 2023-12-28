package main

import (
	"ecommerce/category/db"
	"ecommerce/internal/category"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()
	public := r.Group("/api")

	// category crud
	var categoryRepository = category.NewCategoryRepository(DB)
	var categoryService = category.NewCategoryService(categoryRepository)
	var categoryController = category.NewCategoryController(categoryService)
	public.POST("/category/create", categoryController.CreateCategory)
	public.GET("/category/:id", categoryController.ReadCategory)
	public.PUT("/category/update/:id", categoryController.UpdateCategory)
	public.DELETE("/category/delete/:id", categoryController.DeleteCategory)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

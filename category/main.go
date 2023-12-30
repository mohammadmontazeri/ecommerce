package category

import (
	"ecommerce/category/db"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()
	public := r.Group("/api")

	// category crud
	var categoryRepository = NewCategoryRepository(DB)
	var categoryService = NewCategoryService(categoryRepository)
	var categoryController = NewCategoryController(categoryService)
	public.POST("/category/create", categoryController.CreateCategory)
	public.GET("/category/:id", categoryController.ReadCategory)
	public.PUT("/category/update/:id", categoryController.UpdateCategory)
	public.DELETE("/category/delete/:id", categoryController.DeleteCategory)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

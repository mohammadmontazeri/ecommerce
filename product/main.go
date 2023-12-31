package product

import (
	"ecommerce/product/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()
	public := r.Group("/api")
	// product crud
	var productRepository = NewProductRepository(DB)
	var productService = NewProductService(productRepository)
	var productController = NewProductController(productService)

	public.POST("/product/create", productController.CreateProduct)
	public.GET("/product/:id", productController.ReadProduct)
	public.PUT("/product/update/:id", productController.UpdateProduct)
	public.DELETE("/product/delete/:id", productController.DeleteProduct)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

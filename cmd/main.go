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
	var categoryRepository = category.NewCategoryRepository(DB)
	var categoryService = category.NewCategoryService(categoryRepository)
	var categoryController = category.NewCategoryController(categoryService)
	protected.POST("/category/create", categoryController.CreateCategory)
	protected.GET("/category/:id", categoryController.ReadCategory)
	protected.PUT("/category/update/:id", categoryController.UpdateCategory)
	protected.DELETE("/category/delete/:id", categoryController.DeleteCategory)
	// product crud
	var prodcutWithDB = product.New(DB)

	protected.POST("/product/create", prodcutWithDB.Create)
	protected.GET("/product/:id", prodcutWithDB.Read)
	protected.PUT("/product/update/:id", prodcutWithDB.Update)
	protected.DELETE("/product/delete/:id", prodcutWithDB.Delete)
	// order crud
	var orderRepository = order.NewOrderRepository(DB)
	var orderService = order.NewOrderService(orderRepository)
	var orderController = order.NewOrderController(orderService)
	protected.POST("/order/create", orderController.CreateOrder)
	protected.GET("/order/:id", orderController.ReadOrder)
	protected.PUT("/order/update/:id", orderController.UpdateOrder)
	protected.DELETE("/order/delete/:id", orderController.DeleteOrder)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

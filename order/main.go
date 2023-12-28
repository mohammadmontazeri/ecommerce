package main

import (
	"ecommerce/internal/order"
	"ecommerce/order/db"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()

	public := r.Group("/api")

	// order crud
	var orderRepository = order.NewOrderRepository(DB)
	var orderService = order.NewOrderService(orderRepository)
	var orderController = order.NewOrderController(orderService)
	public.POST("/order/create", orderController.CreateOrder)
	public.GET("/order/:id", orderController.ReadOrder)
	public.PUT("/order/update/:id", orderController.UpdateOrder)
	public.DELETE("/order/delete/:id", orderController.DeleteOrder)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

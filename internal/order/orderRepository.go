package order

import (
	"github.com/gin-gonic/gin"
)

type OrderInterface interface {
	createOrder(c *gin.Context)
}

func NewOrderService(db OrderInterface) *OrderService {
	return &OrderService{db: db}
}

type OrderService struct {
	db OrderInterface
}

func (os *OrderService) Create(c *gin.Context) {
	os.db.createOrder(c)
}

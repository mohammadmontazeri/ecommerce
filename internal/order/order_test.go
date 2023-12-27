package order

import (
	"ecommerce/internal/mocks/ordermocks"
	"ecommerce/models"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOrderWhenOrderWithoutProductNotInsert(t *testing.T) {
	repo := &ordermocks.OrderRepository{}

	repo.On("AddOrderToPivotTable", mock.AnythingOfType("[]int"), mock.AnythingOfType("models.Order")).
		Return(nil).
		Once()

	repo.On("InsertOrderWithoutProducts", mock.AnythingOfType("models.Order")).
		Return(func(db models.Order) models.Order {
			return db
		},
			func(db models.Order) error {
				var err error = errors.New("new error")
				return err
			}).Once()

	orderService := NewOrderService(repo)

	orderInput := models.OrderWithProducts{
		UserID:     1,
		Code:       "ppp",
		Price:      23700,
		Status:     "s",
		ProductsID: []int{1},
	}

	err := orderService.Create(orderInput)

	assert.NotEqual(t, nil, err)
}
func TestCreateOrderWhenProductsIDIsNull(t *testing.T) {
	repo := &ordermocks.OrderRepository{}

	repo.On("AddOrderToPivotTable", mock.AnythingOfType("[]int"), mock.AnythingOfType("models.Order")).
		Return(nil).
		Once()

	repo.On("InsertOrderWithoutProducts", mock.AnythingOfType("models.Order")).
		Return(func(db models.Order) models.Order {
			return db
		}, nil).
		Once()

	orderService := NewOrderService(repo)

	orderInput := models.OrderWithProducts{
		UserID:     1,
		Code:       "ppp",
		Price:      23700,
		Status:     "s",
		ProductsID: []int{},
	}

	err := orderService.Create(orderInput)

	assert.Equal(t, errors.New("error :ProductsID not Corrected !"), err)
}

func TestCreateOrderService(t *testing.T) {

	repo := &ordermocks.OrderRepository{}

	repo.On("AddOrderToPivotTable", mock.AnythingOfType("[]int"), mock.AnythingOfType("models.Order")).
		Return(nil).
		Once()

	repo.On("InsertOrderWithoutProducts", mock.AnythingOfType("models.Order")).
		Return(func(db models.Order) models.Order {
			return db
		}, nil).
		Once()

	orderService := NewOrderService(repo)

	orderInput := models.OrderWithProducts{
		UserID:     1,
		Code:       "ppp",
		Price:      23700,
		Status:     "s",
		ProductsID: []int{1},
	}

	err := orderService.Create(orderInput)

	assert.Nil(t, err)

}

func TestReadOrderService(t *testing.T) {
	repo := &ordermocks.OrderRepository{}

	repo.On("GetOrderFromId", mock.AnythingOfType("int")).
		Return(func(id int) models.Order {
			var db models.Order
			return db
		}, nil).
		Once()

	repo.On("GetOrderProducts", mock.AnythingOfType("models.Order")).
		Return(func(db models.Order) []int {
			var ids []int
			return ids
		}, nil).
		Once()

	orderService := NewOrderService(repo)

	var orderID = 31
	_, err := orderService.Read(orderID)

	assert.Nil(t, err)

}

func TestUpdateOrderService(t *testing.T) {
	repo := &ordermocks.OrderRepository{}

	repo.On("UpdateOrderRow", mock.AnythingOfType("models.Order"), mock.AnythingOfType("int")).
		Return(func(order models.Order, id int) models.Order {
			return order
		}, nil).
		Once()

	repo.On("DeleteOrderProduct", mock.AnythingOfType("int")).
		Return(nil).
		Once()

	repo.On("AddOrderToPivotTable", mock.AnythingOfType("[]int"), mock.AnythingOfType("models.Order")).
		Return(nil).
		Once()

	orderService := NewOrderService(repo)

	orderWithProduct := models.OrderWithProducts{
		UserID:     1,
		Code:       "rrrrrrr",
		Price:      9000,
		Status:     "s",
		ProductsID: []int{1, 3},
	}
	orderID := 30
	err := orderService.Update(orderWithProduct, orderID)

	assert.Nil(t, err)

}

func TestDeleteOrderService(t *testing.T) {
	repo := &ordermocks.OrderRepository{}

	repo.On("DeleteRow", mock.AnythingOfType("models.Order"), mock.AnythingOfType("int")).
		Return(nil).
		Once()

	repo.On("DeleteOrderProduct", mock.AnythingOfType("int")).
		Return(nil).
		Once()

	orderService := NewOrderService(repo)

	orderID := 30
	err := orderService.Delete(orderID)

	assert.Nil(t, err)

}

package order

import (
	ordermocks "ecommerce/internal/mocks/order"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestOrderServiceCreate(t *testing.T) {

	orderMock := new(ordermocks.OrderInterface)

	orderMock.On("createOrder", mock.Anything).Return(nil).Once()

	// om := ordermocks.NewOrderInterface(t)
	// repo := &ordermocks.OrderInterface{}
	// repo := ordermocks.NewOrderInterface(t)

	// repo.On("CreateOrder" , mock.Anything).Return(nil).Once()

	// orderService := NewOrderService(repo)
	// orderService.Create()

}

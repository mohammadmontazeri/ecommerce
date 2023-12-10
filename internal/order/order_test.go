package order

import (
	"bytes"
	"ecommerce/db"
	ordermocks "ecommerce/internal/mocks/order"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func GetTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func MockJsonPost(c *gin.Context, content interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user_id", 1)

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

// func TestCreateOrder(t *testing.T) {

// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request.Method = "POST"
// 	c.Request.Header.Set("Content-Type", "application/json")

// 	var order = []byte(`{	
// 		"id" : 28
// 		"code" : "openorder" ,
// 		"price" : 333,
// 		"user_id" : 1,                            
// 		"products_id" : [1,3] ,
// 		"status" : "a"
// 		}`)

// 	jsonbytes, err := json.Marshal(order)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// the request body must be an io.ReadCloser
// 	// the bytes buffer though doesn't implement io.Closer,
// 	// so you wrap it in a no-op closer
// 	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))



// 	orderMock := ordermocks.NewQueryInterface(t)


// 	orderMock.On("InsertOrderWithoutProducts",order).Return(,nil).Once()

// 	// orderCreate := NewOrderService(orderMock)

// 	// orderCreate.Create(c)
// 	assert.EqualValues(t, http.StatusOK, w.Code)

// 	// orderService.Create()

// }

func TestAddOrderToPivotTable(t *testing.T) {
	
	repo := &ordermocks.QueryInterface{} 

	repo.On("AddOrderToPivotTable",mock.AnythingOfType("[]int"),mock.AnythingOfType("order.Order")).
	Return(nil).
	Once()

	var productsID = []int{1,3}
	order := db.Order{
		UserID: 1,
		Code: "jjjp",
		Price: 3800,
		Status: "s",
	}
	 err := repo.AddOrderToPivotTable(productsID,order)

	 if err !=nil {
		
	 }else{
		
	 }


}
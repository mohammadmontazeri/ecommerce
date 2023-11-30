package order

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	w := httptest.NewRecorder()
	router := gin.Default()
	var jsonStr = []byte(`{	
		"id" : 28
		"code" : "cccccccc" ,
		"price" : 333,
		"user_id" : 1,                            
		"products_id" : [1,3] ,
		"status" : "a"
		}`)
	req, _ := http.NewRequest("POST", "api/admin/order/create", bytes.NewBuffer(jsonStr))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

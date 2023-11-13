package middleware

import (
	"context"
	"ecommerce/configs"
	"ecommerce/internal/product"

	"net/http"

	"github.com/gin-gonic/gin"
)

func VerifyProductCache(c *gin.Context) {
	var prodcut product.Product
	myCache, err := configs.ConnectToRedisForCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	ctx := context.Background()
	id := c.Param("id")

	if err := myCache.Get(ctx, id, &prodcut); err == nil {
		c.JSON(http.StatusOK, gin.H{"product": prodcut})
		c.Abort()
		return
	}

	c.Next()

}

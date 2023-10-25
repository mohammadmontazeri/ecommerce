package middleware

import (
	"fmt"
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtApiMiddleware() gin.HandlerFunc{
	fmt.Println("asdasd")
	return func(c *gin.Context) {
		err := model.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}
		c.Next()
	}
}
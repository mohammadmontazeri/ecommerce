package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtApiMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized User")
			c.Abort()
			return
		}
		c.Next()
	}
}

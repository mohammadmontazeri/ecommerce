package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtApiMiddleware(c *gin.Context) {
	err := TokenValid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UnAuthorized user !"})
		c.Abort()
		return
	}
	c.Next()

	
}

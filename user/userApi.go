package user

import (
	_ "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	u.Email = input.Email
	u.Username = input.Username
	u.Password = input.Password

	// before insert user
	u.BeforeInsert()
	// insert user
	u.Insert(c)

}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}
	u.Username = input.Username
	u.Password = input.Password

	token, err := CheckLogin(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AuthorizedUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "User Authorized !"})
}

package user

import (
	"ecommerce/auth"
	"ecommerce/configs"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

type UserModel struct {
	db *gorm.DB
}

type User struct {
	configs.User
}

func (um *UserModel) Register(c *gin.Context) {
	var input User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}

	u.Email = input.Email
	u.Username = input.Username
	u.Password = input.Password

	// before insert user
	err := u.BeforeInsert()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// insert user
	res := um.db.Create(&u)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User Registration is successful"})

	}

}

func (um *UserModel) Login(c *gin.Context) {
	var input User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := User{}
	u.Username = input.Username
	u.Password = input.Password

	token, err := um.CheckLogin(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func AuthorizedUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Message": "User Authorized !"})
}

func (u *User) BeforeInsert() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func (um *UserModel) CheckLogin(username, password string) (string, error) {

	userId, err := um.CheckUserForLogin(username, password)
	fmt.Println(err)

	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(uint(userId))

	if err != nil {
		return "", err
	}

	return token, nil

}

func (um *UserModel) CheckUserForLogin(username, password string) (int, error) {

	var userLogin User

	res := um.db.Where("username = ?", username).First(&userLogin)

	if res.Error != nil {
		return 0, res.Error
	}

	err := bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, err
	}

	return int(userLogin.ID), nil
}

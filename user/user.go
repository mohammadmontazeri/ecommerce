package user

import (
	"ecommerce/auth"
	"ecommerce/db"
	"fmt"
	"html"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var DB = db.ConnectToDb()

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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

func (u *User) Insert(c *gin.Context) error {

	_, err := DB.Exec("INSERT INTO  users(username,email,password) VALUES($1,$2,$3)", u.Username, u.Email, u.Password)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func CheckLogin(username, password string) (string, error) {

	userId, err := CheckUserForLogin(username, password)
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

func CheckUserForLogin(username, password string) (int, error) {

	var user User
	queryString := fmt.Sprintf("SELECT id,password FROM users WHERE username= '%s'", username)
	err := DB.QueryRow(queryString).Scan(&user.Id, &user.Password)

	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, err
	}

	return user.Id, nil
}

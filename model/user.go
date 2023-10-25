package model

import (
	// "database/sql"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"
)

var db = ConnectToDb()

type User struct {
	ID       int    `json:"id"`
	Username     string `json:"username" binding:"required"`
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

func(u *User) Insert(c *gin.Context)  {

	_, err := db.Exec("INSERT INTO  users(username,email,password) VALUES($1,$2,$3)", u.Username, u.Email, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "User Registration Is Ok"})
	}
}

func CheckLogin(username , password string) (string,error) {
	
	// fmt.Printf("u:%s , i:%s" , User{}.Username ,u)
	err := CheckUserForLogin(username , password)


	if err!=nil {
		return  "" ,err
	}

	token,err := GenerateToken(1)

	if err != nil {
		return "",err
	}

	return token,nil

}

func CheckUserForLogin(username , password string) error {

	var user User 
	queryString := fmt.Sprintf("SELECT password FROM users WHERE username= '%s'",username)
	error := db.QueryRow(queryString).Scan(&user.Password)


	if error != nil {
		return error
	}
	 err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	 if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	 return nil
}
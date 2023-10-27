package auth

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var tokenLifetime = 2 // per hour
var jwtKey = "Montazeeri"

func GenerateToken(user_id uint) (string, error) {

	claim := jwt.MapClaims{}
	claim["user_id"] = user_id
	claim["authorized"] = true
	claim["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifetime)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	return token.SignedString([]byte(jwtKey))

}

func TokenValid(c *gin.Context) error {
	// check token string
	tokenString := c.Query("token")
	if tokenString == "" {
		return errors.New("unAuthorized user !")
	}
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return err
	}
	return nil
}

package user

import (
	"ecommerce/user/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	var DB = db.ConnectToDB()

	r := gin.Default()
	var userWithDB = New(DB)
	public := r.Group("/api")
	public.POST("/register", userWithDB.Register)
	public.POST("/login", userWithDB.Login)

	// protected := r.Group("/api/admin")
	// protected.Use(auth.JwtApiMiddleware)
	// protected.GET("/user", user.AuthorizedUser)

	err := r.Run(":8080")
	if err != nil {
		log.Fatal("server error !")
	}
}

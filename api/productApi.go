package api

import (
	"main/model"

	"github.com/gin-gonic/gin"
)

// create CRUD for Product model

func CreateProduct(c *gin.Context) {
	model.CreateProduct(c)
}

func GetProduct(c *gin.Context) {
	model.GetProduct(c)
}

func UpdateProduct(c *gin.Context) {
	model.UpdateProduct(c)
}

func DeleteProduct(c *gin.Context) {
	model.DeleteProduct(c)
}

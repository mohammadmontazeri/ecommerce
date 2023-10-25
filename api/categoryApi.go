package api

import (
	"main/model"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	model.AddCategory(c)
}

func GetCategory(c *gin.Context){
	model.ReadCategory(c)
}

func UpdateCategory(c *gin.Context){
	model.UpdateCategory(c)
}

func DeleteCategory(c *gin.Context){
	model.DeleteCategory(c)
}
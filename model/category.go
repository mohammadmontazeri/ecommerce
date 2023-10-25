package model

import (
	"fmt"
	"net/http"

	// "reflect"
	"strconv"

	// "strconv"

	"github.com/gin-gonic/gin"
)

type Category struct {
	Id        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Parent_id int    `json:"parent_id"`
}

type AddCategoryInput struct {
	Id        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Parent_id int    `json:"parent_id"`
}

type UpdateCategoryInput struct {
	Name      string `json:"name" binding:"required"`
	Parent_id int    `json:"parent_id"`
}

func AddCategory(c *gin.Context) {
	var input AddCategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := Category{}

	cat.Name = input.Name
	cat.Parent_id = input.Parent_id

	_, err := db.Exec("INSERT INTO  categories(name,parent_id) VALUES($1,$2)", cat.Name, cat.Parent_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Category add successful"})
	}

}

func ReadCategory(c *gin.Context) error {

	var queryParameter = c.Param("id")
	intQueryParameter, _ := strconv.Atoi(queryParameter)

	var cat Category
	if queryParameter == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parmeter not set"})
	}

	queryString := fmt.Sprintf("SELECT * FROM categories WHERE id=%d ", intQueryParameter)
	error := db.QueryRow(queryString).Scan(&cat.Id, &cat.Name, &cat.Parent_id)

	c.JSON(http.StatusOK, gin.H{"category": cat})

	if error != nil {
		return error
	}

	return nil
}

func UpdateCategory(c *gin.Context) {

	var input UpdateCategoryInput
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Model If Exist
	err := GetCategoryIfExisted(c , id)

	if err != nil {
		return
	}

	cat := Category{}
	cat.Name = input.Name
	cat.Parent_id = input.Parent_id

	_, err = db.Exec("UPDATE categories SET name=$1,parent_id=$2 WHERE id=$3", cat.Name, cat.Parent_id, id)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}
}

func DeleteCategory(c *gin.Context) {

	id, error := strconv.Atoi(c.Param("id"))
	if error != nil || id == 0 {
		c.JSON(http.StatusBadRequest,gin.H{"error" :"invalid parameter"})
		return
	}
	// Get Model If Exist
	err := GetCategoryIfExisted(c,id)

	if err != nil {
		return
	}
	_, err = db.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusBadRequest,gin.H{"error" :err.Error()})
	} else {
		c.JSON(http.StatusOK,gin.H{"Message" :"data deleted successful"})	
	}
}

func GetCategoryIfExisted(c *gin.Context , id int) error{
	queryString := fmt.Sprintf("SELECT id FROM categories WHERE id=%d ", id)
	error := db.QueryRow(queryString).Scan(&id)
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "category not found"})
		return error
	}
	return nil
}
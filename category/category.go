package category

import (
	"ecommerce/db"
	"fmt"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type Category struct {
	Id        int    `json:"id"`
	Name      string `json:"name" binding:"required"`
	Parent_id int    `json:"parent_id"`
}

type AddCategoryInput struct {
	Category
}

type UpdateCategoryInput struct {
	Category
}

var DB = db.ConnectToDb()

func Create(c *gin.Context) {
	var input AddCategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := Category{}

	cat.Name = input.Name
	cat.Parent_id = input.Parent_id

	_, err := DB.Exec("INSERT INTO  categories(name,parent_id) VALUES($1,$2)", cat.Name, cat.Parent_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Category add successful"})
	}

}

func Read(c *gin.Context) {

	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var cat Category
	if queryParameter == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parmeter not set"})
		return
	}

	queryString := fmt.Sprintf("SELECT * FROM categories WHERE id=%d ", intQueryParameter)
	err = DB.QueryRow(queryString).Scan(&cat.Id, &cat.Name, &cat.Parent_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": cat})

}

func Update(c *gin.Context) {

	var input UpdateCategoryInput
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := Category{}
	cat.Name = input.Name
	cat.Parent_id = input.Parent_id

	res, err := DB.Exec("UPDATE categories SET name=$1,parent_id=$2 WHERE id=$3", cat.Name, cat.Parent_id, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})

}

func Delete(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := DB.Exec("DELETE FROM categories WHERE id=$1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "data deleted successful"})

}

package category

import (
	"ecommerce/db"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Category struct {
	ID       int    `json:"id"`
	Name     string `binding:"required" json:"name"`
	ParentID int    `json:"parent_id"`
}

type CategoryInput struct {
	Category
}

var DB = db.ConnectToDBGorm()

type Connector interface {
	ConnectDB() *gorm.DB
}

func (cm CategoryModel) ConnectDB() *gorm.DB {
	return db.ConnectToDBGorm()
}
func NewStruct(c Connector) *CategoryModel {
	return &CategoryModel{connector: c}
}

type CategoryModel struct {
	connector Connector
}

func (cm *CategoryModel) Create(c *gin.Context) {
	var input Category

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat := Category{}

	cat.Name = input.Name
	cat.ParentID = input.ParentID

	res := cm.connector.ConnectDB().Create(&cat)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Category add successful"})
	}

}

func (cm *CategoryModel) Read(c *gin.Context) {

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
	res := cm.connector.ConnectDB().Find(&cat, intQueryParameter)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": cat})

}

func (cm *CategoryModel) Update(c *gin.Context) {

	var input CategoryInput
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
	cat.ParentID = input.ParentID

	res := cm.connector.ConnectDB().Model(&cat).Where("id", id).Updates(cat)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})

}

func (cm *CategoryModel) Delete(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := cm.connector.ConnectDB().Delete(&Category{}, id)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "data deleted successful"})

}

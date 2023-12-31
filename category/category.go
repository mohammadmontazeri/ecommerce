package category

import (
	"ecommerce/category/model"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryController model.CategoryService
}

func NewCategoryController(cc model.CategoryService) *CategoryController {
	return &CategoryController{categoryController: cc}
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var input model.Category

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := cc.categoryController.Create(input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Category add successful"})
	}

}

func (cc *CategoryController) ReadCategory(c *gin.Context) {

	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if queryParameter == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parmeter not set"})
		return
	}

	category, err := cc.categoryController.Read(intQueryParameter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	} else {
		c.JSON(http.StatusOK, gin.H{"category": category})
	}

}

func (cc *CategoryController) UpdateCategory(c *gin.Context) {

	var input model.Category
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.categoryController.Update(id, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})

	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}

}

func (cc *CategoryController) DeleteCategory(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.categoryController.Delete(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "data deleted successful"})
	}

}

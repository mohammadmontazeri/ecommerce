package product

import (
	"ecommerce/product/db"
	"ecommerce/product/model"
	"strconv"

	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productController model.ProductService
}

func NewProductController(pc model.ProductService) *ProductController {
	return &ProductController{productController: pc}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {

	var input model.ProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	picturePath, err := UploadProductImage(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = pc.productController.Create(input, picturePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "product add successful"})
	}
}

func (pc *ProductController) ReadProduct(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product model.Product

	rdb, err := db.ConnectToRedis()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	repo := NewRedisRepository(rdb)

	product, err = repo.GetProduct(queryParameter)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"product": product})
	}

	product, err = pc.productController.Read(intQueryParameter)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = repo.SetProduct(queryParameter, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})

}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var input model.ProductInput
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	picturePath, err := UploadProductImage(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = pc.productController.Update(id, input, picturePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "row updated not completed !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}

}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	err = pc.productController.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "row deleted not completed !"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})
	}

}

func UploadProductImage(c *gin.Context) (string, error) {
	file, err := c.FormFile("picture")
	if err != nil {
		return "", err
	}
	filePath := filepath.Join("assets/image", file.Filename)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		return filePath, err
	}
	return filePath, err

}

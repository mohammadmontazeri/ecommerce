package product

import (
	"context"
	"ecommerce/configs"
	"mime/multipart"

	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
)

type ProductInput struct {
	Id         int                   `json:"id"`
	Code       string                `form:"code" json:"code" binding:"required"`
	Title      string                `form:"title" json:"title" binding:"required"`
	Price      float64               `form:"price" json:"price" binding:"required"`
	Picture    *multipart.FileHeader `form:"picture" binding:"required"`
	Detail     string                `form:"detail" json:"detail" binding:"required"`
	CategoryID uint                  `form:"category_id" json:"category_id" binding:"required"`
}

type Product struct {
	configs.Product
}

func New(db *gorm.DB) *ProductModel {
	return &ProductModel{db: db}
}

type ProductModel struct {
	db *gorm.DB
}

func (pm *ProductModel) Create(c *gin.Context) {

	var input ProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	picturePath, err := UploadProductImage(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.CategoryID = input.CategoryID
	product.Picture = picturePath

	res := pm.db.Create(&product)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "product add successful"})
	}
}

func (pm *ProductModel) Read(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product Product

	res := pm.db.Find(&product, intQueryParameter)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}

	myCache, err := configs.ConnectToRedisForCache()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx := context.Background()
	id := c.Param("id")

	err = myCache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   id,
		Value: product,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})

}

func (pm *ProductModel) Update(c *gin.Context) {
	var input ProductInput
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

	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.CategoryID = input.CategoryID
	product.Picture = picturePath

	res := pm.db.Model(&product).Where("id", id).Updates(product)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
}

func (pm *ProductModel) Delete(c *gin.Context) {
	var product Product
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	res := pm.db.Delete(&product, id)

	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.Error.Error()})
		return
	}
	rowsAffected := res.RowsAffected
	if rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no row found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})

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

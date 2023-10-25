package model

import (
	"fmt"
	"log"
	"main/helper"
	"mime/multipart"
	"path/filepath"

	// "mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/gin-gonic/gin/binding"
	// "github.com/gin-gonic/gin/binding"
)

type Image struct {
	Picture *multipart.FileHeader `form:"picture"`
}

type Product struct {
	Id          int     `json:"id"`
	Code        string  `json:"code" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Picture     string  `json:"picture" binding:"required"`
	Detail      string  `json:"detail" binding:"required"`
	Category_id int     `json:"category_id" binding:"required"`
}
type AddProductInput struct {
	// Image
	Code        string                `form:"code" json:"code" binding:"required"`
	Title       string                `form:"title" json:"title" binding:"required"`
	Price       float64               `form:"price" json:"price" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture" binding:"required"`
	Detail      string                `form:"detail" json:"detail" binding:"required"`
	Category_id int                   `form:"category_id" json:"category_id" binding:"required"`
}

type UpdateProductInput struct {
	Code  string  `form:"code" json:"code" binding:"required"`
	Title string  `form:"title" json:"title" binding:"required"`
	Price float64 `form:"price" json:"price" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture" binding:"required"`
	Detail      string `form:"detail" json:"detail" binding:"required"`
	Category_id int    `form:"category_id" json:"category_id" binding:"required"`
}

func CreateProduct(c *gin.Context) {

	var input AddProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// return
	picturePath := UploadProductImage(c)
	// return
	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Picture = picturePath
	product.Detail = input.Detail
	product.Category_id = input.Category_id

	c.JSON(http.StatusOK, gin.H{"message" : product})

	_, err := db.Exec("INSERT INTO  products(code,title,price,picture,detail,category_id) VALUES($1,$2,$3,$4,$5,$6)", product.Code, product.Title, product.Price, product.Picture, product.Detail, product.Category_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "product add successful"})
	}
}

func GetProduct(c *gin.Context) error {
	var queryParameter = c.Param("id")
	intQueryParameter, _ := strconv.Atoi(queryParameter)

	var product Product
	if queryParameter == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parmeter not set"})
	}

	queryString := fmt.Sprintf("SELECT * FROM products WHERE id=%d ", intQueryParameter)
	error := db.QueryRow(queryString).Scan(&product.Id, &product.Title, &product.Code, &product.Price, &product.Picture, &product.Detail, &product.Category_id)

	c.JSON(http.StatusOK, gin.H{"product": product})

	if error != nil {
		return error
	}

	return nil
}

func UpdateProduct(c *gin.Context) {
	var input UpdateProductInput
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Model If Exist
	error := GetModelIfExisted(c, id)

	if error != nil {
		return
	}
	picturePath := UploadProductImage(c)

	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Picture = picturePath
	product.Detail = input.Detail
	product.Category_id = input.Category_id

	_, err := db.Exec("UPDATE products SET title=$1,code=$2,price=$3,detail=$4,category_id=$5,picture=$6 WHERE id=$7", product.Title, product.Code, product.Price, product.Detail, product.Category_id,product.Picture, id)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}
}

func DeleteProduct(c *gin.Context) {
	id, error := strconv.Atoi(c.Param("id"))

	if error != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}
	// Get Model If Exist
	err := helper.GetModelIfExisted("products", id)

	if err != nil {
		return
	}
	_, err = db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "data deleted successful"})
	}
}

// func GetModelIfExisted(c *gin.Context, id int) error {
// 	queryString := fmt.Sprintf("SELECT id FROM products WHERE id=%d ", id)
// 	error := db.QueryRow(queryString).Scan(&id)
// 	if error != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "product not found"})
// 		return error
// 	}
// 	return nil
// }

func UploadProductImage(c *gin.Context) string {
	file, _ := c.FormFile("picture")
	log.Println(file.Filename)
	filePath := filepath.Join("uploads", file.Filename)
	c.SaveUploadedFile(file, filePath)

	return filePath

}

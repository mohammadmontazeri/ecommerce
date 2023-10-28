package product

import (
	"ecommerce/db"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Id          int                   ` json:"id"`
	Code        string                `form:"code" json:"code" binding:"required"`
	Title       string                `form:"title" json:"title" binding:"required"`
	Price       float64               `form:"price" json:"price" binding:"required"`
	Picture     *multipart.FileHeader `form:"picture" binding:"required"`
	Detail      string                `form:"detail" json:"detail" binding:"required"`
	Category_id int                   `form:"category_id" json:"category_id" binding:"required"`
}
type Product struct {
	Id          int     ` json:"id"`
	Code        string  `form:"code" json:"code" binding:"required"`
	Title       string  `form:"title" json:"title" binding:"required"`
	Price       float64 `form:"price" json:"price" binding:"required"`
	Picture     string  `form:"picture" binding:"required"`
	Detail      string  `form:"detail" json:"detail" binding:"required"`
	Category_id int     `form:"category_id" json:"category_id" binding:"required"`
}

var DB = db.ConnectToDb()

func Create(c *gin.Context) {

	var input ProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	picturePath := UploadProductImage(c)
	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.Category_id = input.Category_id

	_, err := DB.Exec("INSERT INTO  products(code,title,price,picture,detail,category_id) VALUES($1,$2,$3,$4,$5,$6)", product.Code, product.Title, product.Price, picturePath, product.Detail, product.Category_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "product add successful"})
	}
}

func Read(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, _ := strconv.Atoi(queryParameter)

	var product Product
	if queryParameter == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "query parmeter not set"})
	}

	queryString := fmt.Sprintf("SELECT * FROM products WHERE id=%d ", intQueryParameter)
	err := DB.QueryRow(queryString).Scan(&product.Id, &product.Title, &product.Code, &product.Price, &product.Picture, &product.Detail, &product.Category_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})

}

func Update(c *gin.Context) {
	var input ProductInput
	id, _ := strconv.Atoi(c.Param("id"))

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	CheckModel(c, id)

	picturePath := UploadProductImage(c)

	product := Product{}

	product.Code = input.Code
	product.Title = input.Title
	product.Price = input.Price
	product.Detail = input.Detail
	product.Category_id = input.Category_id

	_, err := DB.Exec("UPDATE products SET title=$1,code=$2,price=$3,detail=$4,category_id=$5,picture=$6 WHERE id=$7", product.Title, product.Code, product.Price, product.Detail, product.Category_id, picturePath, id)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row updated successful"})
	}
}

func Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	CheckModel(c, id)

	_, err = DB.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})
	}
}

func CheckModel(c *gin.Context, id int) {
	queryString := fmt.Sprintf("SELECT id FROM products WHERE id=%d ", id)
	err := DB.QueryRow(queryString).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "product not found"})
		return
	}
}

func UploadProductImage(c *gin.Context) string {
	file, _ := c.FormFile("picture")
	log.Println(file.Filename)
	filePath := filepath.Join("assets/image", file.Filename)
	c.SaveUploadedFile(file, filePath)

	return filePath

}

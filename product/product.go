package product

import (
	"ecommerce/db"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductInput struct {
	Id         int                   `json:"id"`
	Code       string                `form:"code" json:"code" binding:"required"`
	Title      string                `form:"title" json:"title" binding:"required"`
	Price      float64               `form:"price" json:"price" binding:"required"`
	Picture    *multipart.FileHeader `form:"picture" binding:"required"`
	Detail     string                `form:"detail" json:"detail" binding:"required"`
	CategoryID int                   `form:"category_id" json:"category_id" binding:"required"`
}
type Product struct {
	Id         int     ` json:"id"`
	Code       string  `form:"code" json:"code" binding:"required"`
	Title      string  `form:"title" json:"title" binding:"required"`
	Price      float64 `form:"price" json:"price" binding:"required"`
	Picture    string  `form:"picture" binding:"required"`
	Detail     string  `form:"detail" json:"detail" binding:"required"`
	CategoryID int     `form:"category_id" json:"category_id" binding:"required"`
}

var DB = db.ConnectToDb()

func Create(c *gin.Context) {

	var input ProductInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	picturePath, err := uploadProductImage(c)
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

	_, err = DB.Exec("INSERT INTO  products(code,title,price,picture,detail,category_id) VALUES($1,$2,$3,$4,$5,$6)", product.Code, product.Title, product.Price, picturePath, product.Detail, product.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "product add successful"})
	}
}

func Read(c *gin.Context) {
	var queryParameter = c.Param("id")
	intQueryParameter, err := strconv.Atoi(queryParameter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var product Product

	queryString := fmt.Sprintf("SELECT * FROM products WHERE id=%d ", intQueryParameter)
	err = DB.QueryRow(queryString).Scan(&product.Id, &product.Title, &product.Code, &product.Price, &product.Picture, &product.Detail, &product.CategoryID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})

}

func Update(c *gin.Context) {
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

	picturePath, err := uploadProductImage(c)
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

	res, err := DB.Exec("UPDATE products SET title=$1,code=$2,price=$3,detail=$4,category_id=$5,picture=$6 WHERE id=$7", product.Title, product.Code, product.Price, product.Detail, product.CategoryID, picturePath, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parameter"})
		return
	}

	res, err := DB.Exec("DELETE FROM products WHERE id=$1", id)
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

	c.JSON(http.StatusOK, gin.H{"Message": "row deleted successful"})

}

func uploadProductImage(c *gin.Context) (string, error) {
	file, err := c.FormFile("picture")
	filePath := filepath.Join("assets/image", file.Filename)
	c.SaveUploadedFile(file, filePath)

	return filePath, err

}

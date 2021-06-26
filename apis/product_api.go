package apis

import (
	"fmt"
	"main/database"
	"main/middleware"
	"main/models"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func setupProductAPI(router *gin.Engine) {
	productAPI := router.Group("/api")
	productAPI.Use(middleware.JwtVerify)
	{
		productAPI.GET("/product", getProduct)
		productAPI.GET("/product/:id", getProductByID)
		productAPI.POST("/product", createProduct)
		productAPI.PUT("/product", editProduct)
	}
}

func getProduct(c *gin.Context) {
	var product []models.Product
	keyword := c.Query("keyword")
	if keyword != "" {
		keyword := fmt.Sprintf("%%%s%%", keyword)
		database.GetDB().Where("name Like ?", keyword).Find(&product)
	} else {
		database.GetDB().Find(&product)
	}

	c.JSON(http.StatusOK, product)
}

func getProductByID(c *gin.Context) {
	var product models.Product
	database.GetDB().Where("id = ?", c.Param("id")).Find(&product)
	c.JSON(http.StatusOK, product)
}

func createProduct(c *gin.Context) {
	product := models.Product{}
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	product.CreatedAt = time.Now()
	database.GetDB().Create(&product)
	image, _ := c.FormFile("image")
	saveUploadedFile(image, &product, c)
	c.JSON(http.StatusOK, map[string]interface{}{"status": "create product", "data": product})
}

func saveUploadedFile(image *multipart.FileHeader, product *models.Product, c *gin.Context) {
	if image == nil {
		return
	}
	runningDir, _ := os.Getwd()
	product.Image = image.Filename
	extension := filepath.Ext(image.Filename)
	fileName := fmt.Sprintf("%d%s", product.ID, extension)
	filePath := fmt.Sprintf("%s/uploaded/images/%s", runningDir, fileName)
	if fileExists(filePath) {
		os.Remove(filePath)
	}
	c.SaveUploadedFile(image, filePath)
	database.GetDB().Model(&product).Update("image", fileName)
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsExist(err) {
		return false
	}
	return !info.IsDir()
}

func editProduct(c *gin.Context) {
	var product models.Product
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 64)
	product.ID = uint(id)
	product.Name = c.PostForm("name")
	product.Stock, _ = strconv.ParseInt(c.PostForm("stock"), 10, 64)
	product.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
	database.GetDB().Save(&product)
	image, _ := c.FormFile("image")
	saveUploadedFile(image, &product, c)
	c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": product})
}

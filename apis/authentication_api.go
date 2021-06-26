package apis

import (
	"main/database"
	"main/middleware"
	"main/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthenticationAPI(router *gin.Engine) {
	authenticationAPI := router.Group("/api")
	{
		authenticationAPI.POST("/login", login)
		authenticationAPI.POST("/register", register)

	}
}

func login(c *gin.Context) {
	var user models.User
	if c.ShouldBind(&user) == nil {
		var queryUser models.User
		if err := database.GetDB().First(&queryUser, "username = ?", user.Username).Error; err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{"status": "Unauthorized"})
		} else if checkPasswordHash(user.Password, queryUser.Password) {
			c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": user})
		} else {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{"status": "Unauthorized"})
		}
		user.Password, _ = hashPassword(user.Password)
		user.CreatedAt = time.Now()
		if err := database.GetDB().Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "failed", "error": err})
		} else {
			token := middleware.JwtSign(queryUser)
			c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "token": token})
		}

	} else {
		c.JSON(http.StatusBadRequest, map[string]string{"status": "Bad request"})

	}
}

func register(c *gin.Context) {
	var user models.User
	if c.ShouldBind(&user) == nil {
		user.Password, _ = hashPassword(user.Password)
		user.CreatedAt = time.Now()
		if err := database.GetDB().Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": "failed", "error": err})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{"status": "success", "data": user})
		}

	} else {
		c.JSON(http.StatusBadRequest, map[string]string{"status": "Bad request"})

	}
}

func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

package apis

import (
	"main/database"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	database.SetupDB()
	setupAuthenticationAPI(router)
	setupProductAPI(router)
	setupTransactionAPI(router)
}

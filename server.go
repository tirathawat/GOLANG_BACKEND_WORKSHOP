package main

import (
	"main/apis"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Static("/images", "./uploaded/images")
	apis.Setup(router)
	router.Run(":8081")
}

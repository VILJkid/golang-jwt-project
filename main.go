package main

import (
	"os"

	helper "github.com/VILJkid/golang-jwt-project/helpers"
	"github.com/VILJkid/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	serverPort, isPortSet := os.LookupEnv(helper.SERVER_PORT)
	if !isPortSet {
		serverPort = "8080"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": "Access granted for api-1",
		})
	})

	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"success": "Access granted for api-2",
		})
	})

	router.Run(":" + serverPort)
}

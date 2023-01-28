package main

import (
	"log"
	"os"

	"github.com/VILJkid/golang-jwt-project/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
}

func main() {
	serverPort, isPortSet := os.LookupEnv("SERVER_PORT")
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

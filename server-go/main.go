package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		fmt.Sprintf("http://localhost:%s", os.Getenv("SERVER_PORT")),
		fmt.Sprintf("http://localhost:%s", os.Getenv("CLIENT_MALL_PORT")),
		fmt.Sprintf("http://localhost:%s", os.Getenv("CLIENT_WHADMIN_PORT")),
	}
	router.Use(cors.New(config))

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3003"
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Server is running",
		})
	})

	log.Fatal(router.Run(":" + port))
}

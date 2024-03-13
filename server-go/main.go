package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		fmt.Sprintf("http://localhost:%s", os.Getenv("SERVER_PORT")),
		fmt.Sprintf("http://localhost:%s", os.Getenv("CLIENT_MALL_PORT")),
		fmt.Sprintf("http://localhost:%s", os.Getenv("CLIENT_WHADMIN_PORT")),
	}
	router.Use(cors.New(config))

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	router.Run("localhost:" + port)
	log.Fatal(router.Run("localhost:" + port))
}

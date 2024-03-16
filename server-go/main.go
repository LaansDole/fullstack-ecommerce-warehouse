package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		os.Getenv("DB_USER_ADM"),
		os.Getenv("DB_PASS_ADM"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_DB"),
	))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(gin.Logger())
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

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
		err := db.Ping()
		if err != nil {
			c.JSON(200, gin.H{
				"message": "Only the Server is running",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "Server is connected to MySQL",
			})
		}
	})

	log.Fatal(router.Run(":" + port))
}

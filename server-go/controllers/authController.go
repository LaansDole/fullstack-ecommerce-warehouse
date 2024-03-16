package controllers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	username := c.PostForm("username")
	fmt.Println("Username: ", username) // log username

	password := c.PostForm("password")
	role := c.PostForm("role")
	shop_name := c.PostForm("shop_name")
	city := c.PostForm("city")

	// prevent SQL injection
	matched, _ := regexp.MatchString("^[A-Za-z0-9._-]*$", username)
	if !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The username must not have strange characters"})
		return
	}

	// Check if user exists
	user, err := models.GetLazadaUser(username)
	if err != nil {
		fmt.Println("Error getting Lazada user: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	admin, err := models.GetWHAdmin(username)
	if err != nil {
		fmt.Println("Error getting WHAdmin: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil || admin != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists", "username": username})
		return
	}

	if username == "" || role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a username and role"})
		return
	} else if role == "seller" && shop_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a shop name"})
		return
	}

	if role == "seller" {
		shop, err := models.GetShopName(shop_name)
		if err != nil {
			fmt.Println("Error getting shop name: ", err) // log error
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if shop != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Shop name already exists"})
			return
		}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	fmt.Println(`Hashed Password: ` + string(hashedPassword))

	// Insert the user into the database
	err = models.InsertLazadaUserByRole(role, username, string(hashedPassword), shop_name, city)
	if err != nil {
		fmt.Println("Error inserting user into database: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("User %s created with username: %s", role, username),
		"username":  username,
		"role":      role,
		"shop_name": shop_name,
		"city":      city,
	})
}

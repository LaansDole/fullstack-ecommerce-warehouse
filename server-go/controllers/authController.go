package controllers

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/tokens"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	Role     string `json:"role"`
	ShopName string `json:"shop_name"`
	City     string `json:"city"`
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Username: ", user.Username) // log username

	// prevent SQL injection
	matched, _ := regexp.MatchString("^[A-Za-z0-9._-]*$", user.Username)
	if !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The username must not have strange characters"})
		return
	}

	// Check if user exists
	existingUser, err := models.GetLazadaUser(user.Username)
	if err != nil {
		fmt.Println("Error getting Lazada user: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	admin, err := models.GetWHAdmin(user.Username)
	if err != nil {
		fmt.Println("Error getting WHAdmin: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUser != nil || admin != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists", "username": user.Username})
		return
	}

	if user.Username == "" || user.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a username and role"})
		return
	} else if user.Role == "seller" && user.ShopName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter a shop name"})
		return
	}

	if user.Role == "seller" {
		shop, err := models.GetShopName(user.ShopName)
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	fmt.Println(`Hashed Password: ` + string(hashedPassword))

	// Insert the user into the database
	err = models.InsertLazadaUserByRole(user.Role, user.Username, string(hashedPassword), user.ShopName, user.City)
	if err != nil {
		fmt.Println("Error inserting user into database: ", err) // log error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error inserting user into database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("User %s created with username: %s", user.Role, user.Username),
		"username":  user.Username,
		"role":      user.Role,
		"shop_name": user.ShopName,
		"city":      user.City,
	})
}

func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Username: ", user.Username)
	fmt.Println("Password: ", user.Password)

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide a username and password"})
		return
	}

	var role, shopName string

	// Retrieve the user from the database
	seller, _ := models.GetSeller(user.Username)
	buyer, _ := models.GetBuyer(user.Username)

	if seller != nil {
		role = "seller"
		shopName = seller.ShopName
	} else if buyer != nil {
		role = "buyer"
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	existingUser, _ := models.GetLazadaUser(user.Username)

	fmt.Println("Role: ", role)

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Generate tokens

	userTokens, err := tokens.GenerateTokens(user.Username, role, shopName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Access Tokens: ", userTokens.AccessToken)
	fmt.Println("Refresh Tokens: ", userTokens.RefreshToken)

	// Set the token as a cookie

	tokens.SetTokenCookie(c, user.Username, role, shopName)

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("User %s authenticated", user.Username),
		"username":  user.Username,
		"role":      role,
		"shop_name": shopName,
	})
}

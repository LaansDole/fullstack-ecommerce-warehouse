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
		shopName = ""
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	existingUser, err := models.GetLazadaUser(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	fmt.Println("Role: ", role)

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Generate and set tokens
	generateAndSetTokens(c, user.Username, role, shopName)

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("User %s authenticated", user.Username),
		"username":  user.Username,
		"role":      role,
		"shop_name": shopName,
		"city":      user.City,
	})
}

func LoginAdmin(c *gin.Context) {
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

	existingUser, err := models.GetWHAdmin(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	fmt.Printf("Login user's password_hash: %s\n", existingUser.PasswordHash)

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	role := "admin"
	shopName := ""

	// Generate and set tokens
	generateAndSetTokens(c, user.Username, role, shopName)

	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("User %s authenticated", user.Username),
		"username": user.Username,
		"role":     "admin",
	})
}

func generateAndSetTokens(c *gin.Context, username, role, shopName string) {
	// Generate tokens
	userTokens, err := tokens.GenerateTokens(username, role, shopName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Access Tokens: ", userTokens.AccessToken)
	fmt.Println("Refresh Tokens: ", userTokens.RefreshToken)

	// Set the token as a cookie
	tokens.SetTokenCookie(c, username, role, shopName)
}

func Logout(c *gin.Context) {
	// Get values from the context
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)

	fmt.Printf("%s logged out with role %s\n", username, role)

	// Delete the user's token from the database
	if role == "admin" {
		err := models.DeleteWHAdminToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	} else if role == "seller" || role == "buyer" {
		err := models.DeleteLazadaUserToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Clear the access and refresh tokens
	c.SetCookie("accessToken", "", -1, "/", "", false, true)
	c.SetCookie("refreshToken", "", -1, "/", "", false, true)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message":  fmt.Sprintf("User %s logged out", username),
		"username": username,
		"role":     role,
	})
}

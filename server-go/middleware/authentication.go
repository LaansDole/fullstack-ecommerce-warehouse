package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/tokens"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Authentication() gin.HandlerFunc {
	godotenv.Load()

	return func(c *gin.Context) {
		accessToken, err := c.Cookie("accessToken")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access Token Invalid!"})
			c.Abort()
			return
		}
		refreshToken, err := c.Cookie("refreshToken")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Refresh Token Invalid!"})
			c.Abort()
			return
		}

		fmt.Println("Access token: ", accessToken)
		fmt.Println("Refresh token: ", refreshToken)

		claims, err := tokens.VerifyToken(refreshToken, []byte(os.Getenv("REFRESH_TOKEN_SECRET")))
		if err != nil || claims == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Failed to verify token: " + err.Error()})
			c.Abort()
			return
		}

		username := claims["username"].(string)
		fmt.Println("User: ", username)

		lazadaUser, err := models.GetLazadaUser(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		whAdmin, err := models.GetWHAdmin(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if lazadaUser == nil && whAdmin == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		shopName, ok := claims["shop_name"].(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		tokens.SetTokenCookie(c, username, role, shopName)
		fmt.Printf("%s has a role of %s at Middleware", username, role)

		c.Set("username", username)
		c.Set("role", role)
		c.Set("shop_name", shopName)

		c.Next()
	}
}

package middleware

import (
	"net/http"

	"github.com/LaansDole/fullstack-ecommerce-warehouse/server-go/models"
	"github.com/gin-gonic/gin"
)

func CheckBuyer() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		role := c.MustGet("role").(string)

		buyer, err := models.GetBuyer(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if buyer != nil && role == "buyer" {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: You are not buyer!"})
			c.Abort()
		}
	}
}

func CheckSeller() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		role := c.MustGet("role").(string)

		seller, err := models.GetSeller(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if seller != nil && role == "seller" {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: You are not seller!"})
			c.Abort()
		}
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username").(string)
		role := c.MustGet("role").(string)

		admin, err := models.GetWHAdmin(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if admin != nil && role == "admin" {
			c.Next()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden: You are not admin!"})
			c.Abort()
		}
	}
}

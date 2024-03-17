package middleware

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func FileHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			c.Abort()
			return
		}

		// Ensure the file is an image
		ext := filepath.Ext(file.Filename)
		if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File must be of type JPG, JPEG, or PNG"})
			c.Abort()
			return
		}

		// Save the file to the server
		path := fmt.Sprintf("uploads/%s", file.Filename)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store the file path in the Gin context
		c.Set("imagePath", path)
		fmt.Printf("\nImage path in File Handling Middleware: %s\n", path)

		c.Next()
	}
}

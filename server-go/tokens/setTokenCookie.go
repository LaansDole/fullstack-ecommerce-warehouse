package tokens

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetTokenCookie(c *gin.Context, username string, role string, shopName string) {
	tokens, err := GenerateTokens(username, role, shopName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("access token: ", tokens.AccessToken)
	fmt.Println("refresh token: ", tokens.RefreshToken)

	// Set cookie expiration times
	oneDay := 24 * time.Hour
	longerExp := 30 * 24 * time.Hour

	// Set the access token cookie
	c.SetCookie("accessToken", tokens.AccessToken, int(oneDay.Seconds()), "/", "localhost", false, true)

	// Set the refresh token cookie
	c.SetCookie("refreshToken", tokens.RefreshToken, int(longerExp.Seconds()), "/", "localhost", false, true)

	fmt.Println("response accessToken cookie and refreshToken cookie")
}

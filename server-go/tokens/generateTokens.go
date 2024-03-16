package tokens

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func generateTokens(username, role, shopName string) (Tokens, error) {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var tokens Tokens

	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	if accessTokenSecret == "" || refreshTokenSecret == "" {
		return tokens, fmt.Errorf("token secrets are not set in the environment variables")
	}

	// Generate an access token
	accessToken, err := generateToken(map[string]interface{}{
		"username":  username,
		"role":      role,
		"shop_name": shopName,
		"exp":       time.Now().Add(time.Minute * 30).Unix(),
	}, accessTokenSecret)
	if err != nil {
		return tokens, err
	}

	// Generate a refresh token
	refreshToken, err := generateToken(map[string]interface{}{
		"username":  username,
		"role":      role,
		"shop_name": shopName,
		"exp":       time.Now().AddDate(0, 0, 7).Unix(),
	}, refreshTokenSecret)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken = accessToken
	tokens.RefreshToken = refreshToken

	return tokens, nil
}

func generateToken(claims jwt.MapClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

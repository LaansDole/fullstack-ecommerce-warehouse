package tokens

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func verifyToken(tokenString string, secretKey []byte) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if _, ok := token.Claims.(jwt.StandardClaims); ok && token.Valid {
		return token.Claims, nil
	} else {
		return nil, err
	}
}

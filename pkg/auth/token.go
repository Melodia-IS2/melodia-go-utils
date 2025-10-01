package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ReadClaims(r *http.Request, JWTSecretKey string) (jwt.MapClaims, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, fmt.Errorf("missing authentication token")
	}

	tokenParts := strings.Split(tokenString, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid format for authentication token. Must be 'Bearer myToken'")
	}

	tokenString = tokenParts[1]

	return verifyToken(tokenString, JWTSecretKey)
}

func verifyToken(tokenString string, JWTSecretKey string) (jwt.MapClaims, error) {
	key := []byte(JWTSecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method: %v", token.Method.Alg())
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

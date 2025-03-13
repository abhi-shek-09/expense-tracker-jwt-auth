package utils

import (
	"os"
	"time"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct{
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(id int) (string, error){
	expirationTime := time.Now().Add(24 * time.Hour)
	// make a claims struct with given id, use the expiration time in ur registeredclaims object
	claims := &Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// after the claim is made, sign it with HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (int, error){
	// Extracts the header, payload, and signature from the token.
    // Decodes the payload into the Claims struct.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid{
		return 0, fmt.Errorf("invalid token")
	}
	return claims.UserID, nil
}
package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	ID string `json:"_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(secretKey []byte, id primitive.ObjectID) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := Claims{
		ID: id.Hex(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("ISSUER"),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(secretKey []byte, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("Invalid token or expired")
	}

	// Ensure the token claims are of the correct type
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}

	// Optionally, if you need the ID as ObjectID later on:
	// objectID, err := primitive.ObjectIDFromHex(claims.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("Invalid ObjectID format")
	// }

	return claims, nil
}

func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash from the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

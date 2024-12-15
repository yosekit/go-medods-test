package services

import (
	"crypto/rand"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ExpMins int = 15
)

func GenerateAccessToken(userID string, clientIP string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"ip": clientIP,
		"exp": time.Now().Add(time.Minute * time.Duration(ExpMins)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString(secretKey())

	return accessToken, err
}

func GenerateRefreshToken() []byte {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}

func ParseAccessToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey(), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, errors.New("Invalid access token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid access token claims")
	}

	return claims, nil
}


func secretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

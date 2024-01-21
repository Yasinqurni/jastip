package jwt

import (
	"jastip-app/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
}

type JWTPayload struct {
	SecretKey string
	Expired   float64
	ID        string
}

func GenerateToken(payload JWTPayload) (string, error) {

	expirationTime := time.Now().Add(time.Duration(payload.Expired) * time.Hour)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			ID:        payload.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(payload.SecretKey))
	if err != nil {
		logger.L().Error(err.Error())
		return "", err
	}

	return tokenString, nil
}

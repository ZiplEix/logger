package logger

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ptr[T any](v T) *T {
	return &v
}

func generateToken() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Second * time.Duration(cfg.JwtExpireTime)).Unix(),
	}

	token := jwt.NewWithClaims(cfg.JwtSigningMethod, claims)

	signedToken, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func verifyToken(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, fmt.Errorf("missing or invalid token")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		return false, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); !ok || time.Now().Unix() > int64(exp) {
		return false, fmt.Errorf("token expired")
	}

	return true, nil
}

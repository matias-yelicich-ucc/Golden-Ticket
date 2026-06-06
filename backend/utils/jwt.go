package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the custom claims we include in the JWT
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Rol    string `json:"rol"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	return []byte(secret)
}

// GenerateToken creates a signed JWT with user details
func GenerateToken(userID uint, rol string) (string, error) {
	var expiration time.Duration
	expMinutesStr := os.Getenv("JWT_EXPIRATION_MINUTES")
	if expMinutesStr != "" {
		if minutes, err := strconv.Atoi(expMinutesStr); err == nil {
			expiration = time.Duration(minutes) * time.Minute
		}
	}

	if expiration == 0 {
		expHoursStr := os.Getenv("JWT_EXPIRATION_HOURS")
		expHours := 24
		if expHoursStr != "" {
			if hours, err := strconv.Atoi(expHoursStr); err == nil {
				expHours = hours
			}
		}
		expiration = time.Duration(expHours) * time.Hour
	}

	claims := JWTClaims{
		UserID: userID,
		Rol:    rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateToken parses and checks the validity of a JWT string
func ValidateToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

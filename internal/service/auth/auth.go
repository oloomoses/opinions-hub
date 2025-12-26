package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oloomoses/opinions-hub/internal/dto"
)

var (
	jwtSecretEnvKey      = "JWT_SECRET"
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 14 * time.Hour
)

var jwtSecret = []byte(getEnvOrPanic(jwtSecretEnvKey))

func getEnvOrPanic(key string) string {
	val := os.Getenv(key)

	if val == "" {
		panic("Missing required env: " + key)
	}

	return key
}

func GenerateAccessToken(userID uint, username string) (string, error) {
	claims := dto.JwtClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

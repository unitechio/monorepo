package token

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ResolveBearer(token string) string {
	const prefix = "Bearer "

	if strings.HasPrefix(token, prefix) {
		return strings.TrimPrefix(token, prefix)
	}

	return token
}

func GenerateUUID() string {
	return strings.ReplaceAll(strings.ToLower(jwt.NewNumericDate(time.Now()).String()), ":", "")
}

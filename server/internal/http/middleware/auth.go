package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/unitechio/oss-monorepo/server/pkg/apperr"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

const ContextKeyUserID = "user_id"
const ContextKeyEmail = "email"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			response.FailError(c, apperr.New(http.StatusInternalServerError, "jwt secret is not configured"))
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailError(c, apperr.Unauthorized("authorization header is missing"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			response.FailError(c, apperr.Unauthorized("authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &JWTClaims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			response.FailError(c, apperr.Unauthorized("invalid or expired token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			response.FailError(c, apperr.Unauthorized("invalid token claims"))
			c.Abort()
			return
		}
		if claims.UserID == uuid.Nil {
			response.FailError(c, apperr.Unauthorized("invalid token subject"))
			c.Abort()
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyEmail, claims.Email)
		c.Next()
	}
}

package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/unitechio/oss-monorepo/server/pkg/response"
)

const ContextKeyAccess = "access_profile"

func RequireRoles(authz usecase.AuthorizationService, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := currentUserIDFromContext(c)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		ctx := requestContext(c)
		if err := authz.RequireRoles(ctx, userID, roles...); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		profile, err := authz.GetAccessProfile(ctx, userID)
		if err == nil {
			c.Set(ContextKeyAccess, profile)
		}
		c.Next()
	}
}

func RequireFeature(authz usecase.AuthorizationService, feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := currentUserIDFromContext(c)
		if !ok {
			response.Fail(c, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		if err := authz.RequireFeature(requestContext(c), userID, feature); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.Next()
	}
}

func currentUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextKeyUserID)
	if !exists {
		return uuid.UUID{}, false
	}
	userID, ok := val.(uuid.UUID)
	return userID, ok
}

func requestContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

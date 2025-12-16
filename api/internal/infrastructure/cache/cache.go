package cache

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GetPermissions gets cached permissions for a user
func GetPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// TODO: Implement Redis cache
	return nil, nil
}

// CachePermissions caches permissions for a user
func CachePermissions(ctx context.Context, userID uuid.UUID, permissions []string, ttl time.Duration) error {
	// TODO: Implement Redis cache
	return nil
}

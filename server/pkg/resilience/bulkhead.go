package resilience

import (
	"context"
	"fmt"
)

// ─── Bulkhead (per-server semaphore) ─────────────────────────────────────────

// Bulkhead limits concurrent operations per resource (server).
// This prevents one misbehaving server from consuming all goroutines.
type Bulkhead struct {
	sem  chan struct{}
	name string
}

// NewBulkhead creates a new bulkhead with maxConcurrent slots.
func NewBulkhead(name string, maxConcurrent int) *Bulkhead {
	return &Bulkhead{
		name: name,
		sem:  make(chan struct{}, maxConcurrent),
	}
}

// Execute runs fn within the bulkhead. Returns error if at capacity.
func (b *Bulkhead) Execute(ctx context.Context, fn func() error) error {
	select {
	case b.sem <- struct{}{}:
		defer func() { <-b.sem }()
		return fn()
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("bulkhead %q at capacity — request rejected", b.name)
	}
}

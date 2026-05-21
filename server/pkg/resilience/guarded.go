package resilience

import (
	"context"
	"errors"
)

// ─── Combined: Retry + CB guard ──────────────────────────────────────────────

// GuardedCall wraps fn with circuit breaker + retry + backoff.
// Ideal for: dispatcher → agent command send.
func GuardedCall(
	ctx context.Context,
	cb *CircuitBreaker,
	retryCfg RetryConfig,
	fn func() error,
) error {
	return ExponentialBackoff(ctx, retryCfg, func(err error) bool {
		// Don't retry circuit-open errors — they'll just fail immediately anyway
		return !errors.Is(err, ErrCircuitOpen)
	}, func() error {
		return cb.Execute(fn)
	})
}

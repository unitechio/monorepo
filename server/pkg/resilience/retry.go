package resilience

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"

	"github.com/rs/zerolog/log"
)

// ─── Retry with Exponential Backoff ──────────────────────────────────────────

// RetryConfig defines retry behaviour.
type RetryConfig struct {
	MaxAttempts int           // max total attempts (including first)
	BaseDelay   time.Duration // initial delay before first retry
	MaxDelay    time.Duration // cap on delay
	Multiplier  float64       // backoff multiplier (typically 2.0)
	Jitter      float64       // jitter fraction [0.0, 1.0] — prevents thundering herd
}

// DefaultRetryConfig is tuned for agent command dispatch.
var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 4,
	BaseDelay:   250 * time.Millisecond,
	MaxDelay:    10 * time.Second,
	Multiplier:  2.0,
	Jitter:      0.2,
}

// IsRetryable lets callers classify errors as retryable or permanent.
type IsRetryable func(err error) bool

// ExponentialBackoff retries fn with exponential backoff + jitter.
// Returns the last error if all attempts fail.
// If ctx is cancelled, returns ctx.Err() immediately.
func ExponentialBackoff(ctx context.Context, cfg RetryConfig, isRetryable IsRetryable, fn func() error) error {
	var lastErr error
	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if isRetryable != nil && !isRetryable(lastErr) {
			return lastErr // permanent error — don't retry
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		delay := delay(cfg, attempt)
		log.Warn().
			Err(lastErr).
			Int("attempt", attempt).
			Dur("backoff", delay).
			Msg("[retry] retrying after backoff")

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
	return fmt.Errorf("all %d attempts failed: %w", cfg.MaxAttempts, lastErr)
}

func delay(cfg RetryConfig, attempt int) time.Duration {
	base := float64(cfg.BaseDelay) * math.Pow(cfg.Multiplier, float64(attempt-1))
	if base > float64(cfg.MaxDelay) {
		base = float64(cfg.MaxDelay)
	}
	if cfg.Jitter > 0 {
		base = base * (1 + cfg.Jitter*(rand.Float64()*2-1))
	}
	if base < 0 {
		base = 0
	}
	return time.Duration(base)
}

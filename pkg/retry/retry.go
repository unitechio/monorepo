package retry

import (
	"math/rand/v2"
	"time"

	"github.com/rs/zerolog/log"
)

type Settings struct {
	MaxRetries    int
	RetryTimeFunc func(failure int) time.Duration
}

var DefaultSettings = Settings{
	MaxRetries:    3,
	RetryTimeFunc: RetryTime,
}

func RetryWithWarning[T any](operation string, settings Settings, fn func() (T, error)) (T, error) {
	timeStart := time.Now()

	var result T
	var err error

	for attempt := 0; attempt <= settings.MaxRetries; attempt++ {
		result, err = fn()
		if err == nil {
			if attempt > 1 {
				log.Info().
					Str("operation", operation).
					Int("attempt", attempt).
					Dur("duration", time.Since(timeStart)).
					Msg("operation completed successfully")
			}
			return result, nil
		}

		if attempt < settings.MaxRetries {
			retryDelay := settings.RetryTimeFunc(attempt)
			log.Warn().
				Err(err).
				Str("operation", operation).
				Int("attempt", attempt).
				Int("max_retries", settings.MaxRetries).
				Str("retry_delay", retryDelay.String()).
				Str("duration", time.Since(timeStart).String()).
				Msg("operation failed, retrying")
			time.Sleep(retryDelay)
		}
	}

	log.Error().
		Err(err).
		Str("operation", operation).
		Int("attempt", settings.MaxRetries).
		Dur("duration", time.Since(timeStart)).
		Msg("operation failed after max retries")
	return result, err
}

func RetryTime(failure int) time.Duration {
	var baseTime time.Duration

	switch failure {
	case 0, 1:
		baseTime = 5 * time.Second
	case 2:
		baseTime = 30 * time.Second
	default:
		baseTime = 90 * time.Second
	}

	//plus random jitter
	jitter := time.Duration(rand.Int64N(int64(10 * time.Second)))
	return baseTime + jitter
}

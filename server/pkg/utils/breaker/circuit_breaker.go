package breaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	CLOSED State = iota
	OPEN
	HALF_OPEN
)

type CircuitBreaker struct {
	mu sync.Mutex

	failures      int
	successes     int
	state         State
	threshold     int
	timeout       time.Duration
	lastFailTime  time.Time
	halfOpenLimit int
}

func New(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:         CLOSED,
		threshold:     threshold,
		timeout:       timeout,
		halfOpenLimit: 3,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case OPEN:
		if time.Since(cb.lastFailTime) > cb.timeout {
			cb.state = HALF_OPEN
			cb.successes = 0
			return true
		}
		return false
	default:
		return true
	}
}

func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.state = CLOSED
}

func (cb *CircuitBreaker) Fail() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.threshold {
		cb.state = OPEN
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return errors.New("circuit open")
	}

	err := fn()
	if err != nil {
		cb.Fail()
		return err
	}

	cb.Success()
	return nil
}

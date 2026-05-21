package resilience

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// ─── Circuit Breaker ──────────────────────────────────────────────────────────

// CBState is the circuit breaker state.
type CBState int32

const (
	CBClosed   CBState = 0 // healthy — requests pass through
	CBOpen     CBState = 1 // failing — requests rejected immediately
	CBHalfOpen CBState = 2 // probing — one request allowed through
)

func (s CBState) String() string {
	switch s {
	case CBClosed:
		return "closed"
	case CBOpen:
		return "open"
	case CBHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// ErrCircuitOpen is returned when the circuit is open.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker protects a downstream dependency.
type CircuitBreaker struct {
	name string

	// config
	failureThreshold int           // failures before opening
	successThreshold int           // successes in half-open before closing
	resetTimeout     time.Duration // open → half-open transition

	// state (all atomic / mutex-protected)
	state         atomic.Int32
	failures      int
	successes     int
	lastFailureAt time.Time
	mu            sync.Mutex

	// optional metrics callback: state, name
	onStateChange func(name string, from, to CBState)
}

// CBConfig configures the circuit breaker.
type CBConfig struct {
	Name             string
	FailureThreshold int
	SuccessThreshold int
	ResetTimeout     time.Duration
	OnStateChange    func(name string, from, to CBState)
}

// NewCircuitBreaker creates a new CircuitBreaker.
func NewCircuitBreaker(cfg CBConfig) *CircuitBreaker {
	cb := &CircuitBreaker{
		name:             cfg.Name,
		failureThreshold: cfg.FailureThreshold,
		successThreshold: cfg.SuccessThreshold,
		resetTimeout:     cfg.ResetTimeout,
		onStateChange:    cfg.OnStateChange,
	}
	// start closed
	cb.state.Store(int32(CBClosed))
	return cb
}

// DefaultAgentCB creates a circuit breaker tuned for agent command dispatch.
func DefaultAgentCB(serverID string, onStateChange func(string, CBState, CBState)) *CircuitBreaker {
	return NewCircuitBreaker(CBConfig{
		Name:             "agent-" + serverID,
		FailureThreshold: 5,
		SuccessThreshold: 2,
		ResetTimeout:     30 * time.Second,
		OnStateChange:    onStateChange,
	})
}

// Execute runs fn through the circuit breaker.
// Returns ErrCircuitOpen if the circuit is open.
func (cb *CircuitBreaker) Execute(fn func() error) error {
	state := CBState(cb.state.Load())

	switch state {
	case CBOpen:
		cb.mu.Lock()
		sinceFailure := time.Since(cb.lastFailureAt)
		cb.mu.Unlock()
		if sinceFailure < cb.resetTimeout {
			return ErrCircuitOpen
		}
		// Transition to half-open: allow one probe
		cb.transitionTo(CBHalfOpen)

	case CBHalfOpen:
		// Only one concurrent probe request
	}

	err := fn()

	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.recordFailure()
	} else {
		cb.recordSuccess()
	}
	return err
}

func (cb *CircuitBreaker) recordFailure() {
	current := CBState(cb.state.Load())
	cb.failures++
	cb.lastFailureAt = time.Now()
	cb.successes = 0

	if current == CBHalfOpen || cb.failures >= cb.failureThreshold {
		cb.transitionTo(CBOpen)
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	current := CBState(cb.state.Load())
	cb.successes++

	if current == CBHalfOpen && cb.successes >= cb.successThreshold {
		cb.failures = 0
		cb.transitionTo(CBClosed)
	}
}

func (cb *CircuitBreaker) transitionTo(next CBState) {
	prev := CBState(cb.state.Swap(int32(next)))
	if prev == next {
		return
	}
	log.Info().
		Str("cb", cb.name).
		Str("from", prev.String()).
		Str("to", next.String()).
		Msg("[circuit-breaker] state changed")
	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, next)
	}
}

// State returns the current circuit state.
func (cb *CircuitBreaker) State() CBState { return CBState(cb.state.Load()) }

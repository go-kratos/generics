package retry

import (
	"math/rand"
	"time"
)

// backoffConfig stores the exponential backoff parameters.
type backoffConfig struct {
	baseDelay time.Duration
	maxDelay  time.Duration
	mult      float64
	jitter    float64
}

func defaultBackoff() backoffConfig {
	return backoffConfig{
		baseDelay: 100 * time.Millisecond,
		maxDelay:  15 * time.Second,
		mult:      1.6,
		jitter:    0.2,
	}
}

func (bc backoffConfig) duration(retries int) time.Duration {
	if retries == 0 {
		return bc.baseDelay
	}
	backoff, max := float64(bc.baseDelay), float64(bc.maxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.mult
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so callers that fail together avoid lockstep retries.
	backoff *= 1 + bc.jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}

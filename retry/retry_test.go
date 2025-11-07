package retry

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetrySucceedsAfterTransientFailures(t *testing.T) {
	var attempts int32
	r := New(
		WithAttempts(3),
		WithBaseDelay(time.Microsecond),
		WithMaxDelay(time.Microsecond),
		WithMultiplier(1),
		WithJitter(0),
	)

	err := r.Do(context.Background(), func(context.Context) error {
		if atomic.AddInt32(&attempts, 1) < 3 {
			return errors.New("temporary")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("retry returned unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryStopsWhenNotRetryable(t *testing.T) {
	wantErr := errors.New("permanent")
	var calls int
	r := New(
		WithAttempts(5),
		WithRetryable(func(err error) bool { return !errors.Is(err, wantErr) }),
	)

	gotErr := r.Do(context.Background(), func(context.Context) error {
		calls++
		return wantErr
	})
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, gotErr)
	}
	if calls != 1 {
		t.Fatalf("expected 1 attempt, got %d", calls)
	}
}

func TestBackoffOptionsApplied(t *testing.T) {
	base := 10 * time.Millisecond
	max := time.Second
	mult := 2.5
	jitter := 0.15
	r := New(
		WithBaseDelay(base),
		WithMaxDelay(max),
		WithMultiplier(mult),
		WithJitter(jitter),
	)

	if r.backoff.baseDelay != base {
		t.Fatalf("expected base delay %v, got %v", base, r.backoff.baseDelay)
	}
	if r.backoff.maxDelay != max {
		t.Fatalf("expected max delay %v, got %v", max, r.backoff.maxDelay)
	}
	if r.backoff.mult != mult {
		t.Fatalf("expected multiplier %v, got %v", mult, r.backoff.mult)
	}
	if r.backoff.jitter != jitter {
		t.Fatalf("expected jitter %v, got %v", jitter, r.backoff.jitter)
	}
}
